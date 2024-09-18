
const uname = prompt('input your nickname');
const ws = new WebSocket("ws://" + document.location.host + "/ws");
const status    = document.getElementById("status");
const msg_box   = document.getElementById("msg_box");
const msg_list  = document.getElementById("msg_list");
const user_list = document.getElementById("user_list");
const user_num  = document.getElementById("user_num");

ws.onopen = function () {
   status.innerHTML = "Successfully connected to server";
   var msg = {
      "action": "user-join",
      "sender": uname
   };
   send(msg)
};

ws.onmessage = function (e) {
   var msg = JSON.parse(e.data);

   switch (msg.action) {
      case 'send-message':
         appendMsg(msg.sender + ': ' + msg.content);
         break;
      case 'user-join':
         appendMsg(msg.sender + " jointed")
         list(msg.user_list)
         break;
      case 'user-left':
         appendMsg(msg.sender + " left")
         list(msg.user_list)
         break;
   }
};

ws.onerror = function () {
   status.innerHTML = "Connection error";
};

ws.onclose = function () {
   status.innerHTML = "Connection closed";
}

function confirm(event) {
   var key_num = event.keyCode;
   if (13 == key_num) {
      sendMsg();
   } else {
      return false;
   }
}

function send(msg) {
   var data = JSON.stringify(msg);
   ws.send(data);
}

function sendMsg() {
   var content = msg_box.value;
   var reg = new RegExp("\r\n", "g");
   content = content.replace(reg, "").trim();
   msg_box.value = '';

   if (content == "") {
      return
   }
   var msg = {
      "content": content,
      "action": "send-message",
      "sender": uname
   };
   send(msg);
}

function appendMsg(data) {
   var msg = document.createElement("p");
   msg.innerHTML = data;
   msg_list.appendChild(msg);
   msg_list.scrollTop = msg_list.scrollHeight;
}

function list(list){
   if (list == null) 
      return
   
   while (user_list.hasChildNodes()) {
      user_list.removeChild(user_list.firstChild);
   }
   for (var index in list) {
      var user = document.createElement("p");
      user.innerHTML = list[index];
      user_list.appendChild(user);
   }
   user_num.innerHTML = list.length;
   user_list.scrollTop = user_list.scrollHeight;
}