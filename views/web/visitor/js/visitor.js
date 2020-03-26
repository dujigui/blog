var mainNav = document.getElementById("js-menu");
var navBarToggle = document.getElementById("js-navbar-toggle");

navBarToggle.addEventListener("click", function() {
    mainNav.classList.toggle("active");
});

var s = document.createElement("script");
s.type = "text/javascript";
s.charset = "UTF-8";

var host = window.location.hostname;
if (host.includes("dujigui")) {
    s.src = "http://tajs.qq.com/stats?sId=66502163";
} else if (host.includes("nullsfootprints")) {
    s.src = "http://tajs.qq.com/stats?sId=66502158";
}

document.getElementById("tencent_analysis").append(s);