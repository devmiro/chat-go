const socket = new WebSocket(`ws://${window.location.host}/chatsocket/{{.RoomName}}`);
const messageInput = document.getElementById('message-input');
const sendButton = document.getElementById('send-button');
const chatBox = document.getElementById('chat-box');

socket.onmessage = function (event) {
    const message = event.data;
    const p = document.createElement('p');
    p.textContent = message;
    chatBox.appendChild(p);
};

sendButton.addEventListener('click', function () {
    const message = messageInput.value;
    if (message.trim() !== '') {
        socket.send(message);
        messageInput.value = '';
    }
});
