fetch('/api/docs')
  .then(res => res.json())
  .then(docs => {
    const sidebar = document.getElementById('sidebar');
    const details = document.getElementById('details');

    // Grupowanie endpointów po kategorii
    const grouped = {};
    docs.forEach(doc => {
      if (!grouped[doc.category]) grouped[doc.category] = [];
      grouped[doc.category].push(doc);
    });

    // Dla każdej kategorii stwórz kontener z możliwością rozwijania
    Object.entries(grouped).forEach(([category, endpoints]) => {
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
        const btn = document.createElement('button');
        btn.className = `endpoint-btn`;

       const hasLuaFunc = !!doc.luaFunc;

      btn.innerHTML = `
          <span class="method-badge method-${doc.method.toLowerCase()}">${doc.method}</span>
          <span class="endpoint-text">${doc.endpoint}</span>
          ${hasLuaFunc ? `<span class="test-badge">test</span>` : ''}
      `;

        btn.onclick = () => {
          details.innerHTML = `
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

          // req sim part
          details.innerHTML += `
            <h3>Request Simulator</h3>
            <div id="simulator">
              <label for="req-method">Method:</label><br>
              <select id="req-method">
                <option>GET</option>
                <option>POST</option>
                <option>PUT</option>
                <option>DELETE</option>
              </select><br><br>

              <label for="req-headers">Headers (JSON):</label><br>
              <textarea id="req-headers" rows="4" style="width:100%; font-family: monospace;">${doc.headers}</textarea><br><br>

              <label for="req-body">Body:</label><br>
              <textarea id="req-body" rows="6" style="width:100%; font-family: monospace;">${doc.body}</textarea><br><br>

              <button id="send-request">Send Request</button>
              <div id="simulator-response" style="margin-top: 1em;"></div>
            </div>
          `;

          document.getElementById('send-request').onclick = async () => {
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

            resBox.innerHTML = "<em>Sending request to Lua...</em>";

            try {
              const res = await fetch("/api/simulate", {
                method: "POST",
                headers: {
                  "Content-Type": "application/json"
                },
                body: JSON.stringify({
                  endpoint: doc.endpoint,
                  method,
                  headers,
                  body
                })
              });

              const json = await res.json();
              console.log(json.response[0])
              resBox.innerHTML = `
                <strong>Status:</strong> ${json.response[0]}<br>
                <strong>Response:</strong><pre>${JSON.stringify(json.response[1], null, 2)}</pre>
                <strong>Log:</strong><pre>${(json.log || []).join('\n')}</pre>
              `;

            } catch (err) {
              resBox.innerHTML = `<span style="color:red;"><strong>Error:</strong> ${err}</span>`;
            }
          };

        };

        list.appendChild(btn);
      });


      section.appendChild(header);
      section.appendChild(list);
      sidebar.appendChild(section);
    });
  });