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
            document.title = bn + ' - Home';
            break;
        case 'search':
            eScript.src = '/visitor/js/search.js';
            document.title = bn + ' - Search';
            break;
        case 'about':
            document.title = bn + ' - About Me';
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