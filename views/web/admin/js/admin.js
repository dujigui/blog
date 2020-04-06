layui.use(['element', 'layer'], function () {
    $ = layui.jquery;
    var layer = layui.layer;
    toast = function (msg) {
        layer.msg(msg);
    };

    var s = document.createElement('script');
    var tn = '';
    switch (tab) {
        case 'dashboard':
            $('#dashboard').addClass("layui-this");
            s.src = "/backyard/js/dashboard.js";
            tn = '概览';
            break;
        case 'post_list':
            $('#post').addClass("layui-this");
            $('#post_list').addClass("layui-this");
            s.src = "/backyard/js/post_list.js";
            tn = '文章列表';
            break;
        case 'post_editor':
            $('#post').addClass("layui-this");
            $('#post_editor').addClass("layui-this");
            s.src = "/backyard/js/post_editor.js";
            tn = '文章编辑器';
            break;
        case 'file':
            $('#file').addClass("layui-this");
            s.src = "/backyard/js/file.js";
            tn = '文件';
            break;
        case 'tag':
            $('#tag').addClass("layui-this");
            s.src = "/backyard/js/tag.js";
            tn = '标签';
            break;
        case 'comment':
            $('#comment').addClass("layui-this");
            s.src = "/backyard/js/comment.js";
            tn = '评论';
            break;
        case 'preference':
            $('#preference').addClass("layui-this");
            s.src = "/backyard/js/preference.js";
            tn = '设置';
            break;
    }
    document.title = apn + '-' + tn;
    document.body.appendChild(s);

    $('#btnLogout').click(function () {
        $.ajax({
            type: "DELETE",
            contentType: "application/json",
            url: "/logout",
            context: document.body,
            success: function (result) {
                toast("已退出登录");
                setTimeout(function () {
                    window.location.href = "/login"
                }, 1000);
            },
            error: function (result) {
                console.log(result);
            }
        });
    });
    $.get("/admin/info", function (result) {
        $('#avatar').attr('src', result.result.avatar);
        $('#nickname').text(result.result.nickname);
    })
});

function formatCreated(d) {
    return formatDate(d.Created)
}

function formatUpdated(d) {
    return formatDate(d.Updated)
}

function formatDate(date) {
    var d = new Date(date);
    return `${d.getFullYear()}-${("0" + (d.getMonth() + 1)).slice(-2)}-${("0" + d.getDate()).slice(-2)} ${("0" + d.getHours()).slice(-2)}:${("0" + d.getMinutes()).slice(-2)}`
}