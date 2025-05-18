class ApiDocsUI {
  constructor(sidebarId = 'sidebar', detailsId = 'details') {
    this.sidebar = document.getElementById(sidebarId);
    this.details = document.getElementById(detailsId);
    this.docs = [];
    this.grouped = {};
  }

  static init() {
    const ui = new ApiDocsUI();
    ui.fetchDocs();
    return ui;
  }

  async fetchDocs() {
    try {
      const res = await fetch(`/api/docs?_=${Date.now()}`);
      this.docs = await res.json();
      this.groupDocs();
      this.renderSidebar();
    } catch (err) {
      this.sidebar.innerHTML = `<div style="color:red;">Błąd ładowania dokumentacji: ${err}</div>`;
    }
  }

  groupDocs() {
    this.grouped = {};
    this.docs.forEach(doc => {
      if (!this.grouped[doc.category]) this.grouped[doc.category] = [];
      this.grouped[doc.category].push(doc);
    });
  }

  renderSidebar() {
    this.sidebar.innerHTML = ''; // wyczyść przed renderowaniem
    Object.entries(this.grouped).forEach(([category, endpoints]) => {
      const section = document.createElement('div');
      const header = document.createElement('button');
      header.textContent = category.toUpperCase();
      header.className = 'category-btn';

      const list = document.createElement('div');
      list.style.display = 'none';

      header.onclick = () => {
        list.style.display = list.style.display === 'none' ? 'block' : 'none';
      };

      endpoints.forEach(doc => {
        const container = this.createEndpointContainer(doc);
        list.appendChild(container);
      });

      section.appendChild(header);
      section.appendChild(list);
      this.sidebar.appendChild(section);
    });
  }

  createEndpointContainer(doc) {
    const btn = document.createElement('button');
    btn.className = 'endpoint-btn';

    const hasLuaFunc = !!doc.luaFunc;
    const hasTestTable = !!doc.defaultDB;
    const hadMarkdown = !!doc.markdown;

    btn.innerHTML = `
      <span class="method-badge method-${doc.method.toLowerCase()}">${doc.method}</span>
      <span class="endpoint-text">${doc.endpoint}</span>
      ${hasLuaFunc ? `<span class="test-badge">test</span>` : ''}
      ${hasTestTable ? `<span class="table-badge">DB</span>` : ''}
      ${hadMarkdown ? `<span class="markdown-badge">MD</span>` : ''}
    `;

    const container = document.createElement('div');
    container.className = "endpoint-container";
    container.appendChild(btn);

    // markdown-y POD endpointem
    if (doc.markdown && Array.isArray(doc.markdown)) {
      const markdownWrapper = document.createElement('div');
      markdownWrapper.className = "markdown-wrapper";
      doc.markdown.forEach(path => {
        const mdBtn = document.createElement('button');
        mdBtn.textContent = path.split('/').pop();
        mdBtn.className = 'markdown-btn';

        mdBtn.onclick = async () => {
          const res = await fetch(`/api/markdowns/view?path=${encodeURIComponent(path)}`);
          const mdText = await res.text();
          const html = marked.parse(mdText);
          this.details.innerHTML = `<h2>${path}</h2>${html}`;
        };

        markdownWrapper.appendChild(mdBtn);
      });
      container.appendChild(markdownWrapper);
    }

    btn.onclick = () => this.showEndpointDetails(doc);
    return container;
  }

  showEndpointDetails(doc) {
    this.details.innerHTML = `
      <h2>${doc.method} ${doc.endpoint}</h2>
      <p><em>${doc.description}</em></p>
      <div><strong>Permissions:</strong> ${doc.permissions}</div>
      <div><strong>Request Body:</strong><pre>${doc.body}</pre></div>
      <div><strong>Request Headers:</strong><pre>${doc.headers}</pre></div>
      <div><strong>Response:</strong><pre>${doc.res}</pre></div>
      <div><strong>Errors:</strong><ul>
        ${doc.errors.map(e => `<li><code>${e.code}</code>: ${e.message} – ${e.description}</li>`).join('')}
      </ul></div>
    `;

    // Testowa tabela
    if (doc.defaultDB && doc.defaultDB.length > 0) {
      const headers = Object.keys(doc.defaultDB[0]);
      const tableHTML = `
        <h3>Test Database (DefaultDB)</h3>
        <table class="test-db-table">
          <thead>
            <tr>
              ${headers.map(h => `<th>${h}</th>`).join('')}
            </tr>
          </thead>
          <tbody>
            ${doc.defaultDB.map(row => `
              <tr>
                ${headers.map(h => `<td>${row[h]}</td>`).join('')}
              </tr>
            `).join('')}
          </tbody>
        </table>
      `;
      this.details.innerHTML += tableHTML;
    }

    this.details.innerHTML += this.getSimulatorHTML(doc);

    setTimeout(() => {
      this.setupSimulator(doc);
    }, 0); // po renderze
  }

  getSimulatorHTML(doc) {
    const default_method = doc.method.toUpperCase();
    return `
      <h3>Request Simulator</h3>
      <div id="simulator">
        <label for="req-method">Method:</label><br>
        <select id="req-method">
          <option value="GET" ${default_method === "GET" ? "selected" : ""}>GET</option>
          <option value="POST" ${default_method === "POST" ? "selected" : ""}>POST</option>
          <option value="PUT" ${default_method === "PUT" ? "selected" : ""}>PUT</option>
          <option value="DELETE" ${default_method === "DELETE" ? "selected" : ""}>DELETE</option>
        </select><br><br>
        <label for="req-headers">Headers (JSON):</label><br>
        <textarea id="req-headers" rows="4" style="width:100%; font-family: monospace;">${doc.headers}</textarea><br><br>
        <label for="req-body">Body:</label><br>
        <textarea id="req-body" rows="6" style="width:100%; font-family: monospace;">${doc.body}</textarea><br><br>
        <button id="send-request">Send Request</button>
        <div id="simulator-response" style="margin-top: 1em;"></div>
      </div>
    `;
  }

  setupSimulator(doc) {
    const btn = document.getElementById('send-request');
    if (!btn) return;
    btn.onclick = async () => {
      const method = document.getElementById('req-method').value;
      const headersRaw = document.getElementById('req-headers').value;
      const body = document.getElementById('req-body').value;
      const resBox = document.getElementById('simulator-response');

      let headers = {};
      try {
        headers = JSON.parse(headersRaw);
      } catch (e) {
        resBox.innerHTML = "<span style='color:red;'>Invalid headers JSON</span>";
        return;
      }

      const dbToSend = doc.defaultDB && Array.isArray(doc.defaultDB)
        ? doc.defaultDB.map(entry => ({ ...entry }))
        : [];

      resBox.innerHTML = "<em>Sending request to Lua...</em>";

      try {
        const res = await fetch("/api/simulate", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            endpoint: doc.endpoint,
            method,
            headers,
            body,
            defaultDB: dbToSend
          })
        });

        const json = await res.json();
        const logArray = Array.isArray(json.log)
          ? json.log
          : Object.values(json.log || {});

        resBox.innerHTML = `
          <strong>Status:</strong> ${json.response?.status || 'N/A'}<br>
          <strong>Response:</strong><pre>${typeof json.response.body === 'string' ? json.response.body : JSON.stringify(json.response.body, null, 2)}</pre>
          <strong>Log:</strong><pre>${logArray.join('\n')}</pre>
        `;

        if (json.db) {
          const dbArray = Array.isArray(json.db)
            ? json.db
            : Object.values(json.db);

          if (dbArray.length === 0) return;

          const headersSet = new Set();
          dbArray.forEach(row => {
            if (typeof row === 'object') {
              Object.keys(row).forEach(k => headersSet.add(k));
            }
          });
          const headers = Array.from(headersSet);

          const tableHTML = `
            <h4>Updated DB:</h4>
            <table class="test-db-table">
              <thead>
                <tr>${headers.map(h => `<th>${h}</th>`).join('')}</tr>
              </thead>
              <tbody>
                ${dbArray.map(row => `
                  <tr>
                    ${headers.map(h => `<td>${row[h] ?? ''}</td>`).join('')}
                  </tr>
                `).join('')}
              </tbody>
            </table>
          `;

          resBox.innerHTML += tableHTML;
        }

      } catch (err) {
        resBox.innerHTML = `<span style="color:red;"><strong>Error:</strong> ${err}</span>`;
      }
    };
  }
}

document.addEventListener('DOMContentLoaded', () => {
  ApiDocsUI.init();
});

// ====================================

class MarkdownsUI {
  constructor(sidebarId = 'sidebar', detailsId = 'details') {
    this.sidebar = document.getElementById(sidebarId);
    this.details = document.getElementById(detailsId);
    this.mdFiles = [];
    this.groupedMd = {};
  }

  static init() {
    const ui = new MarkdownsUI();
    ui.fetchMarkdowns();
    return ui;
  }

  async fetchMarkdowns() {
    try {
      const res = await fetch(`/api/markdowns?_=${Date.now()}`);
      this.mdFiles = await res.json();
      this.groupMarkdowns();
      this.renderSidebar();
    } catch (err) {
      this.sidebar.innerHTML = `<div style="color:red;">Błąd ładowania markdownów: ${err}</div>`;
    }
  }

  groupMarkdowns() {
    this.groupedMd = {};
    this.mdFiles.forEach(file => {
      if (!this.groupedMd[file.category]) this.groupedMd[file.category] = [];
      this.groupedMd[file.category].push(file);
    });
  }

  renderSidebar() {
    Object.entries(this.groupedMd).forEach(([category, files]) => {
      const visibleFiles = files.filter(file => !file.name.startsWith('_'));
      if (visibleFiles.length === 0) return;

      const section = document.createElement('div');
      const header = document.createElement('button');
      header.textContent = `[MD] ${category.toUpperCase()}`;
      header.className = 'category-btn';

      const list = document.createElement('div');
      list.style.display = 'none';

      header.onclick = () => {
        list.style.display = list.style.display === 'none' ? 'block' : 'none';
      };

      visibleFiles.forEach(file => {
        const btn = this.createMdBtn(file);
        list.appendChild(btn);
      });

      section.appendChild(header);
      section.appendChild(list);
      this.sidebar.appendChild(section);
    });
  }

  createMdBtn(file) {
    const btn = document.createElement('button');
    btn.className = `endpoint-btn method-md`;
    btn.innerHTML = `
      <span class="method-badge method-md">.md</span>
      <span class="endpoint-text">${file.name}</span>
    `;

    btn.onclick = () => this.showMarkdown(file);
    return btn;
  }

  async showMarkdown(file) {
    try {
      const res = await fetch(`/api/markdowns/view?path=${encodeURIComponent(file.path)}`);
      const mdText = await res.text();

      const html = marked.parse(mdText);
      this.details.innerHTML = `<h2>${file.name}</h2>${html}`;
    } catch (err) {
      this.details.innerHTML = `<span style="color:red;">Błąd ładowania pliku: ${err}</span>`;
    }
  }
}

document.addEventListener('DOMContentLoaded', () => {
  MarkdownsUI.init();
});
