// Server Info Client

let latestInterfaceInfo = [];
let selectedInterfaceName = "";

// Utility functions
const getCurrentPort = () => {
  const url = new URL(window.location.href);
  return url.port ? url.port : '80';
};

const createTable = (array, tableId) => {
  const table = document.getElementById(tableId);
  table.innerHTML = '';

  array.forEach(device => {
    const row = document.createElement('tr');
    const cell = document.createElement('td');

    cell.innerHTML = device;
    row.appendChild(cell);
    table.appendChild(row);
  });
};

const clickableInterfaceTable = (interfaces, tableId) => {
  const table = document.getElementById(tableId);
  table.innerHTML = '';
  document.getElementById('port').innerText = `Server Port: ${getCurrentPort()}`;

  interfaces.forEach((iface) => {
    const clickableInterface = document.createElement('a');

    clickableInterface.textContent = `${iface.ip} (${iface.name})`;
    clickableInterface.className = 'pure-button';
    clickableInterface.style.margin = '4px';
    clickableInterface.href = '#';

    clickableInterface.addEventListener('click', (event) => {
      event.preventDefault();

      if (iface.name != selectedInterfaceName) {
        document.querySelector("#qr").src="data:image/png;base64," + iface.qr;
        selectedInterfaceName = iface.name;
      }
    });

    table.append(clickableInterface)
  });
};

document.addEventListener("DOMContentLoaded", () => {
  const pairingElement = document.getElementById('pairingStatus');
  const pairingDescriptionElement = document.getElementById('pairingDescription');
  const pairingInfoDiv = document.getElementById('pairingInfo');
  const pairingLoadingDiv = document.getElementById('pairingLoading');
  const deviceSectionDiv = document.getElementById('deviceSection');
  const secretText = document.getElementById('secret');

  console.log("Hello, ActionPad!");
  pairingInfoDiv.style.display = 'none';

  setInterval(() => {
    // TODO: Make sure you use the same port as currently running server instance

    // Make GET request to status endpoint
    fetch('/status')
    .then(response => {
      deviceSection.style.display = 'block';

      if (!response.ok) {
        throw new Error(`An error has occurred: ${response.status}`);
      }

      return response.json();
    })
    .then(data => {
      const { connected, saved } = data;

      if (connected) {
        createTable(connected, 'connectedTable');
      } else {
        document.getElementById('connectedTable').innerHTML = '<p>No devices connected.</p>';
        console.log('No devices are connected');
      }

      if (saved && saved.length > 0) {
        createTable(saved, 'savedTable');
      } else {
        document.getElementById('savedTable').innerHTML = '<p>No devices saved.</p>';
        console.log('No saved devices');
      }
    })
    .catch(error => {
      console.error(`Fetch error: ${error}`);
    });

    fetch('/interfaces')
    .then(response => {
      if (!response.ok) {
        throw new Error(`An error has occurred: ${response.status}`);
      }

      return response.json();
    })
    .then(interfaces => {
      if (Array.isArray(interfaces)) {
        latestInterfaceInfo = interfaces

        if (selectedInterfaceName.length == 0 && interfaces.length > 0) {
          const iface = interfaces[0];
          document.querySelector("#qr").src="data:image/png;base64," + iface.qr;
          selectedInterfaceName = iface.name;
        }

        if (interfaces.length == 0) {
          document.getElementById('ipTable').innerHTML = 'No IP addresses available.';
        }

        clickableInterfaceTable(interfaces, 'ipTable')
      } else {
        throw new Error(`Interfaces response is not an array.`);
      }
    })
    .catch(error => {
      console.error(`Fetch error: ${error}`);
    });

    // Make GET request to pairing endpoint
    fetch('/pairing')
    .then(response => {
      pairingLoadingDiv.style.display = 'none';

      if (!response.ok) {
        pairingElement.style.color = '#fc3d39';
        pairingElement.innerText = 'Pairing mode disabled.';
        pairingInfoDiv.style.display = 'none';
        
        pairingDescriptionElement.innerText = 'New devices are NOT allowed to connect. You must enable pairing mode from the ActionPad Server menu in order for new devices to connect.';

        throw new Error(`An error has occurred: ${response.status}`);
      }
      

      return response.json();
    })
    .then(data => {
      const { name, code } = data;
      pairingElement.style.color = '#53d769';
      pairingElement.innerText = 'Pairing mode enabled.';
      pairingInfoDiv.style.display = 'block';
      secretText.innerText = 'Server Code: ' + code;
      pairingDescriptionElement.innerText = 'New devices are allowed to connect. Once you have finished connecting your device(s), disable pairing mode from the server menu for additional security.';
    });
  }, 1000);

  function parse_query_string(e){for(var n=e.split("&"),o={},r=0;r<n.length;r++){var t=n[r].split("="),s=decodeURIComponent(t[0]),d=decodeURIComponent(t[1]);if(void 0===o[s])o[s]=decodeURIComponent(d);else if("string"==typeof o[s]){var p=[o[s],decodeURIComponent(d)];o[s]=p}else o[s].push(decodeURIComponent(d))}return o}var parsed_qs=parse_query_string(window.location.href);
  document.querySelector('#secret').textContent = 'Server Code: ' + parsed_qs[Object.keys(parsed_qs)[0]];
  document.querySelector("#qr").src="data:image/png;base64,"+parsed_qs[Object.keys(parsed_qs)[1]];
});

