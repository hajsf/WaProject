let sseUri = "http://localhost:1234/sse"
let source = new EventSource(sseUri);
let reconnecting = false;
let live = document.querySelector('.badge')
// live.style.backgroundColor = '#ff334b';
/*
function connect() {
  source = new EventSource("http://localhost:1234/sse");
}
connect();
setInterval(() => {
    if (source.readyState == EventSource.CLOSED) {
        reconnecting = true;
        console.log("reconnecting...");
        connect();
    } else if (reconnecting) {
        reconnecting = false
        console.log("reconnected!");
    }
}, 3000); */

source.onopen = function(){
  live = document.querySelector('.badge')
  live.style.backgroundColor = '#ff334b';
}

source.onerror = function() {
  live.style.backgroundColor = '#30000614';
  setInterval(() => {
    if (source.readyState == EventSource.CLOSED) {
        live.style.backgroundColor = '#30000614';
        source.close();
        source = new EventSource("http://localhost:1234/sse");
    } 
}, 3000);
}

source.onmessage = function (event) {
  console.log(event.data)
}

// will be called automatically whenever the server sends a message with the event field set to "qr"
// echo "event: notification\ndata: {"time": "' . $curDate . '"}';
// fmt.Fprint(w, "event: notification\ndata: %v\n\n", c) then followed by fmt.Fprintf(w, "data: %v\n\n", c)

source.addEventListener("notification", function(event) {
  console.log(event.data)
  document.querySelector('#qr').innerHTML = "";
  var message = event.data
  document.querySelector('#message').innerHTML = message;
});

source.addEventListener("qrCode", function(event) {
  console.log(event.data)
  document.querySelector('#qr').innerHTML = "";
  document.querySelector('#message').innerHTML = "Scan the QR code from the WhatsApp application";
  var qrcode = new QRCode("qr", {
    text: message,
    width: 128,
    height: 128,
    colorDark : "#000000",
    colorLight : "#ffffff",
    correctLevel : QRCode.CorrectLevel.M
});
});