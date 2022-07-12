let sseUri = "http://localhost:1234/sse/signal"
let source = new EventSource(sseUri);
let reconnecting = false;
let live = document.querySelector('.badge')
live.style.backgroundColor = '#ff334b';

function connect() {
  source = new EventSource("http://localhost:1234/sse/signal");
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
}, 3000);

source.onopen = function(){
  live = document.querySelector('.badge')
  live.style.backgroundColor = '#ff334b';
}

source.onerror = function() {
  live.style.backgroundColor = '#30000614';
  source.close();
}

source.onmessage = function (event) {

}

// will be called automatically whenever the server sends a message with the event field set to "qr"
// echo "event: qr\ndata: {"time": "' . $curDate . '"}';
// fmt.Fprint(w, "event: qr\ndata: %v\n\n", c) then followed by fmt.Fprintf(w, "data: %v\n\n", c)
source.addEventListener("qr", function(event) {

});

