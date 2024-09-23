
const uname = prompt('input your nickname');
const ws = new WebSocket("ws://" + document.location.host + "/ws");
const status    = document.getElementById("status");
const mymsg     = document.getElementById("mymsg");
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

   if (msg.type == "image"){
      var name = document.createElement("p");
      name.innerHTML = msg.sender + ":";
      msg_list.appendChild(name);

      var img = document.createElement('img');
      img.src = msg.content;
      img.style.display = "block";
      msg_list.appendChild(img);

      msg_list.scrollTop = msg_list.scrollHeight;
   } else {
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
   var content = mymsg.value;
   var reg = new RegExp("\r\n", "g");
   content = content.replace(reg, "").trim();
   mymsg.value = '';

   if (content == "") {
      return
   }
   var msg = {
      "action": "send-message",
      "content": content,
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

function upload() {
   var fileInput = document.getElementById('fileInput');
   if (fileInput.files.length == 0) {
      return
   }

   var file = fileInput.files[0];

   var formData = new FormData();
   formData.append("myFile", file);

   var request = new XMLHttpRequest();
   request.open("POST", "/upload");

   request.onload = () => {
      fileInput.value = '';
      var response = JSON.parse(request.responseText);
      var file = response.filename;
      console.log(file)
      var msg = {
         "action": "upload-image",
         "content": file,
         "sender": uname
      };
      send(msg)
   };

  request.send(formData);

}