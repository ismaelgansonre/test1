const personTable = document.getElementById('personTable').getElementsByTagName('tbody')[0];

function addOrUpdatePersonInTable(person) {
    // Check if the person already exists in the table
    let existingRow = Array.from(personTable.rows).find(row => row.cells[0].textContent === person.name);

    if (existingRow) {
        // Update the existing row
        existingRow.cells[1].textContent = person.age;
    } else {
        // Add a new row
        const row = personTable.insertRow();
        const nameCell = row.insertCell(0);
        const ageCell = row.insertCell(1);

        nameCell.textContent = person.name;
        ageCell.textContent = person.age;
    }
}

document.getElementById('createPersonForm').addEventListener('submit', function(event) {
    event.preventDefault();
    const name = document.getElementById('name').value;
    const age = document.getElementById('age').value;
    fetch('/create_person', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ name: name, age: parseInt(age) })
    })
    .then(response => {
        if (!response.ok) {
            if (response.status === 409) {
                alert('Person already exists.');
            } else {
                throw new Error('Network response was not ok ' + response.statusText);
            }
            return null;
        }
        return response.json();
    })
    .then(data => {
        if (data) {
            alert('Person created: ' + JSON.stringify(data));
            addOrUpdatePersonInTable(data); // Add or update the person in the table
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
});

document.getElementById('updatePersonForm').addEventListener('submit', function(event) {
    event.preventDefault();
    const name = document.getElementById('updateName').value;
    const age = document.getElementById('updateAge').value;
    fetch('/update_person_age', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({ name: name, age: parseInt(age) })
    })
    .then(response => {
        if (!response.ok) {
            if (response.status === 404) {
                alert('Person not found.');
            } else {
                throw new Error('Network response was not ok ' + response.statusText);
            }
            return null;
        }
        return response.json();
    })
    .then(data => {
        if (data) {
            alert('Person updated: ' + JSON.stringify(data));
            addOrUpdatePersonInTable(data); // Add or update the person in the table
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
});

const fetchProtectedDataButton = document.getElementById('fetchProtectedData');
const authModal = document.getElementById('authModal');
const authForm = document.getElementById('authForm');
const closeSpan = document.getElementsByClassName('close')[0];

fetchProtectedDataButton.onclick = function() {
    authModal.style.display = 'block';
};

closeSpan.onclick = function() {
    authModal.style.display = 'none';
};

window.onclick = function(event) {
    if (event.target == authModal) {
        authModal.style.display = 'none';
    }
};

authForm.onsubmit = function(event) {
    event.preventDefault();
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    fetch('/protected', {
        headers: {
            'Authorization': 'Basic ' + btoa(username + ':' + password)
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }
        return response.json();
    })
    .then(data => {
        authModal.style.display = 'none';
        document.getElementById('protectedData').textContent = JSON.stringify(data, null, 2);
    })
    .catch(error => {
        alert('Authentication failed.');
        console.error('Error:', error);
    });
};
