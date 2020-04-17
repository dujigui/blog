layui.use(['form'], function () {
    var $ = layui.jquery;
    var form = layui.form;
    $('#qq').click(function () {
        window.location.href=`https://graph.qq.com/oauth2.0/authorize?response_type=code&client_id=${qqAppId}&redirect_uri=${encodeURI(redirectUri)}&state=${qqState}`;
    });

    form.on('submit(formDemo)', function (data) {
        if (data.field.account === '') {
            layer.msg('账号不能为空');
            return false
        }

        if (data.field.password === '') {
            layer.msg('密码不能为空');
            return false
        }

        var loginUrl = window.location.pathname + window.location.search;


        $.ajax({
            type: "POST",
            contentType: "application/json",
            url: loginUrl,
            context: document.body,
            data: JSON.stringify(data.field),
            dataType: 'json',
            success: function (result) {
                console.log("result: " + result);
                layer.msg('登录成功');
                setTimeout(function () {
                    var urlParams = new URLSearchParams(window.location.search);
                    var redirect = decodeURI(urlParams.get('redirect') || "/");
                    window.location.href = redirect
                }, 1000);
            },
            error: function (result) {
                console.error(result);
                layer.msg('登录失败: ' + result.responseJSON.msg);
            }
        });

        return false;
    });
});