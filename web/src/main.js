var streamEndPoint = "http://localhost:3000/stream"
var filterEndPoint = "http://localhost:3000/filter"

var filterVal = ''

var source = new EventSource(streamEndPoint);
source.onmessage = function(event) {startSse(event)};

function startSse(event){
    // document.getElementById("result").innerHTML += event.data + "<br>";
    // console.log("got: " + event.data)
    if (filterVal.len == 0) {
        $("#log-box[readonly]").html($("#log-box[readonly]").val() + event.data + "\n");
    } else {
        var match = false
        var filterValSplit = filterVal.split("|")
        eventSplit = event.data.split("\n")
        for (var x=0; x<eventSplit.length; x++) {
            for (var i=0; i<filterValSplit.length; i++) {
                // console.log("checking " + filterValSplit[i])
                if (eventSplit[x].indexOf(filterValSplit[i].trim()) != -1) {
                    // console.log("match!")
                    match = true
                } else {
                    match = false
                    break
                }
            }
            if (match) {
                // console.log("matched: " + eventSplit[x])
                $("#log-box[readonly]").html($("#log-box[readonly]").val() + eventSplit[x] + "\n");
            }
        }
    }
    // console.log("val of box: " + $("#log-box[readonly]").val())
    var $textarea = $("#log-box[readonly]");
    $textarea.scrollTop($textarea[0].scrollHeight);
    // $("#log-box[readonly]").scrollTop($("#log-box[readonly]").scrollHeight;
}

$("#stream-btn").click(function(event) {
    console.log("TOGGLE STREAM")
    if (source == null) {
        console.log("Enable STREAM")
        $("#stream-btn").val("Streaming...")
        $("#stream-btn").addClass("success")
        $("#stream-btn").removeClass("warning")
        source = new EventSource(streamEndPoint);
        source.onmessage = function(event){startSse(event)}
    } else {
        console.log("Disable STREAM")
        $("#stream-btn").val("Paused")
        $("#stream-btn").addClass("warning")
        $("#stream-btn").removeClass("success")
        source.onmessage = null;
        source.close();
        source = null;
    }
})

$("#filter-box").keyup(function(event){
    if(event.keyCode == 13){
        $("#filter-btn").click();
    }
});

$("#filter-btn").click(function(event) {
    filterVal = $("#filter-box").val();
    // if (filterVal.length > 0) {
        console.log("Filtering: " + filterVal)
        $.get(filterEndPoint + "?q=" + filterVal,
        function(data, status){
            if (status == "success"){
                console.log("Got data")
                $("#log-box[readonly]").html(data)
            }
            // console.log("status: " + status)
            // console.log("Data Loaded: " + data);
        });
    // } else {
    //     console.log("nothign to filter")
    // }
})
