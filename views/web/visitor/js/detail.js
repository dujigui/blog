layui.use('form', function(){
    var form = layui.form;

    form.on('submit(formDemo)', function(data){
        $.ajax({
            type: "POST",
            contentType: "application/json",
            url: "/comments",
            context: document.body,
            data: JSON.stringify(data.field),
            dataType: 'json',
            success: function (result) {
                toast("评论成功");
                setTimeout(function () {
                    window.location.reload();
                }, 1000);
            },
            error: function (result) {
                console.log(result);
                toast("评论失败");
            }
        });
        return false;
    });
});

$('#btnLogin').click(function () {
    window.location.href = "/login?redirect=" + encodeURI(window.location.pathname)
});

$('#btnLogout').click(function () {
    $.ajax({
        type: "DELETE",
        contentType: "application/json",
        url: "/logout",
        context: document.body,
        success: function (result) {
            toast("已退出登录");
            setTimeout(function () {
                window.location.reload();
            }, 1000);
        },
        error: function (result) {
            console.log(result);
        }
    });
});