layui.use('element', function () {
    var element = layui.element;

    var now = new Date();
    var year = now.getFullYear();
    $('#year-pc-h').append(year + " 使用进度")

    // 过了多少天
    var start = new Date(now.getFullYear(), 0, 0);
    var diff = (now - start) + ((start.getTimezoneOffset() - now.getTimezoneOffset()) * 60 * 1000);
    var od = 1000 * 60 * 60 * 24;
    var day = Math.floor(diff / od);

    // 今年多少天
    var total;
    if (year % 400 === 0 || (year % 100 !== 0 && year % 4 === 0)) {
        total = 366;
    } else {
        total = 365;
    }

    setTimeout(function () {
        element.progress('year-pc', ((day/total) * 100).toFixed(2) + '%');
    }, 500);
});