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
        btn.textContent = `${doc.method} ${doc.endpoint}`;
        btn.className = 'endpoint-btn';
        btn.onclick = () => {
          details.innerHTML = `
            <h2>${doc.method} ${doc.endpoint}</h2>
            <p><em>${doc.description}</em></p>
            <div><strong>Permissions:</strong> ${doc.permissions}</div>
            <div><strong>Request Body:</strong><pre>${doc.body}</pre></div>
            <div><strong>Response:</strong><pre>${doc.res}</pre></div>
            <div><strong>Errors:</strong><ul>
              ${doc.errors.map(e => `<li><code>${e.code}</code>: ${e.message} – ${e.description}</li>`).join('')}
            </ul></div>
          `;
        };
        list.appendChild(btn);
      });

      section.appendChild(header);
      section.appendChild(list);
      sidebar.appendChild(section);
    });
  });
