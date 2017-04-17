var streamEndPoint = "http://localhost:3000/stream"
var filterEndPoint = "http://localhost:3000/filter"

var filterVal = ''
var source

var $textarea = $("#log-box[readonly]");
$.get(filterEndPoint + "?q=" + filterVal,
    function(data, status){
        if (status == "success"){
            // console.log("Got data")
            $("#log-box[readonly]").html(data)
            source = new EventSource(streamEndPoint);
            source.onmessage = function(event) {startSse(event)};
            $textarea.scrollTop($textarea[0].scrollHeight);
        }
        // console.log("status: " + status)
        // console.log("Data Loaded: " + data);
    }
);

function startSse(event){
    var autoScroll = false
    var originalScrollTop = $textarea.scrollTop()
    if (originalScrollTop + ($textarea[0].scrollHeight/20) >= $textarea[0].scrollHeight) {
        autoScroll = true
    }

    var currentVal = $("#log-box[readonly]").val()
    if (filterVal == '') {
        $("#log-box[readonly]").html(currentVal + event.data + "\n");
        // console.log("Trimming data length: ", currentVal.length)
        trimLargeTextarea(currentVal)
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
                $("#log-box[readonly]").html(currentVal + eventSplit[x] + "\n");
                trimLargeTextarea(currentVal)
            }
        }
    }
    // console.log("val of box: " + $("#log-box[readonly]").val())
    if (autoScroll) {
        $textarea.scrollTop($textarea[0].scrollHeight);
    } else {
        $textarea.scrollTop(originalScrollTop)
    }
    // $("#log-box[readonly]").scrollTop($("#log-box[readonly]").scrollHeight;
}

$("#stream-btn").click(function(event) {
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
                // console.log("Got data")
                data = data.split("\n").reverse().join("\n") + "\n";
                $("#log-box[readonly]").html(data)
                $textarea.scrollTop($textarea[0].scrollHeight);
            }
            // console.log("status: " + status)
            // console.log("Data Loaded: " + data);
        });
    // } else {
    //     console.log("nothign to filter")
    // }
})

$(function() {
  $("#filter-box").focus();
});

function trimLargeTextarea(currentVal){
    if (currentVal.length > 1200000) {
        console.log("Trimming data length: ", currentVal.length)
        $("#log-box[readonly]").html(currentVal.substr(currentVal.length/2*-1));
    }
}