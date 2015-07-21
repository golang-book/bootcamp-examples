function init() {
  // the tweet overlay
  var tweetOverlay = document.getElementById("tweet-overlay");
  var tweetOverlayForm = document.getElementById("tweet-overlay-form");
  tweetOverlayForm.addEventListener("submit", function(evt) {
    evt.preventDefault();

    var text = document.getElementById("tweet-overlay-text").value;
    var xhr = new XMLHttpRequest();
    xhr.open("POST", "/api/tweets");
    xhr.send(JSON.stringify({
      Text: text
    }));
    xhr.onreadystatechange = function() {
      if (xhr.readyState === 4) {
        if (xhr.status === 200) {
          location.reload(false);
        } else {
          alert(xhr.responseText);
        }
      }
    };
  });
  // hide the overlay after clicking cancel
  var tweetOverlayCancel = document.getElementById("tweet-overlay-cancel");
  tweetOverlayCancel.addEventListener("click", function(evt) {
    evt.preventDefault();

    tweetOverlay.style.display = "none";
  });
  // hide the overlay on hitting ESC
  document.body.addEventListener("keyup", function(evt) {
    if (evt.keyCode === 27) {
      tweetOverlay.style.display = "none";
    }
  });

  // the tweet button on the top right of the page
  var tweetButton = document.getElementById("tweet-button");
  tweetButton.addEventListener("click", function(evt) {
    evt.preventDefault();

    tweetOverlay.style.display = "block";
  });
}

init();
