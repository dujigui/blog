layui.use(['element', 'layer'], function () {
    $ = layui.jquery;
    var layer = layui.layer;
    toast = function (msg) {
        layer.msg(msg);
    };

    $('#js-navbar-toggle').click(function () {
        $('#js-menu').toggleClass("active");
    });

    var eScript = document.createElement('script');
    switch (tab) {
        case 'home':
            document.title = bn + 'Home';
            break;
        case 'search':
            eScript.src = '/visitor/js/search.js';
            document.title = bn + 'Search';
            break;
        case 'about':
            document.title = bn + 'About Me';
            break;
        case 'detail':
            eScript.src = '/visitor/js/detail.js';
            document.title = bn;
            break;
        default:
            document.title = bn;
    }
    document.body.appendChild(eScript);
});

//todo 更好地整合腾讯统计
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