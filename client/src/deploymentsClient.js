
function deploymentsClient() {
    
    var ws = new WebSocket("ws://localhost:8080/ws");
    ws.onmessage = (event) => {        
        this.onmessage(JSON.parse(event.data))       
    }
}

deploymentsClient.prototype.onMessage  = function(onmessage) {
    this.onmessage = onmessage;
};

export default deploymentsClient;
