'use strict';

ws.onclose = (evt) => {
    console.log("Closeed websocket connection");
    ws = null;
}

ws.onerror = function(evt) {
    console.log("Websocket error: " + evt.data);
}

ws.onopen = (evt) => {
    console.log("Opened websocket connection");
}

function Service(id, name) {
    this.id = id;
    this.name = name;
}

function sortServices(a, b) {
    if (a.name > b.name) {
        return 1;
    }

    if (a.name < b.name) {
        return -1;
    }

    return 0;
}

let links = [];
document.querySelectorAll('#containerList>a').forEach((link) => {
    let name = link.innerHTML.trim();
    if (name.includes('-')) {
        name = name.substr(0, name.lastIndexOf('-')).trim();
    }
    links.push(new Service(link.id, name));
});

ws.onmessage = function(evt) {
    const data = JSON.parse(evt.data);
    if (!('Action' in data)) {
        console.log('Invalid json response: ' + evt.data);
        return;
    }

    console.log(data);

    switch (data['Action']) {
        case 'start':
            addEntry(data);
            break;
        case 'die':
            removeEntry(data);
            break;
        case 'err':
            if (!('ErrMessage' in data)) {
                console.log('Received error from socket connection, but no message was sent');
                break;
            }
            console.log(data['ErrMessage']);
            break;
    }
}

function addEntry(data) {
    if (data['Action'] !== "start") {
        console.log('Invalid action for addEntry: ' + data['Action']);
        return;
    }

    if (!('Container' in data)) {
        console.log('Object passed in to addEntry does not have Container field');
        return;
    }

    if (!('Info' in data)) {
        console.log('Object passed in to addEntry does not have Info field');
        return;
    }

    if (!('Name' in data['Info'])) {
        console.log('The Info object does not contain a Name field');
        return;
    }

    if (!('Description' in data['Info'])) {
        console.log('The Info object does not contain a Description field');
        return;
    }

    if (!('Url' in data['Info'])) {
        console.log('The Info object does not contain a Url field');
        return;
    }

    if (document.getElementById(data['Container']) != null) {
        console.log('Trying to add duplicate entry');
        return;
    }

    const containerList = document.getElementById('containerList');
    if (containerList == null) {
        console.log('Missing target for links');
        return;
    }

    const serviceName = data['Info']['Name'];
    const description = data['Info']['Description'];
    const url = data['Info']['Url'];
    let linkEntry = '';

    if (description.trim() !== "") {
        linkEntry = `
        <a id="${data['Container']}" href="${url}" class="list-group-item list-group-item-action">
            ${serviceName} - ${description}
        </a>`
    } else {
        linkEntry = `
        <a id="${data['Container']}" href="${url}" class="list-group-item list-group-item-action">
            ${serviceName}
        </a>`
    }

    const s = new Service(data['Container'], serviceName);
    links.push(s);
    links.sort(sortServices);

    const position = links.indexOf(s);
    if ((links.length === 1) || (position + 1 === links.length)) {
        containerList.innerHTML += linkEntry;
        return;
    }

    const template = document.createElement('template');
    linkEntry = linkEntry.trim();
    template.innerHTML = linkEntry;
    containerList.insertBefore(template.content.firstChild, document.getElementById(links[position + 1].id));
}

function removeEntry(data) {
    if (data['Action'] !== "die") {
        console.log('Invalid action for addEntry: ' + data['Action']);
        return;
    }

    if (!('Container' in data)) {
        console.log('Object passed in to removeEntry does not have Container field');
        return;
    }

    links = links.filter(services => services.id !== data['Container']);

    const entry = document.getElementById(data['Container']);
    if (entry != null) {
        entry.remove();
    }
}