layui.use('form', function () {
    var form = layui.form;
    var layer = layui.layer;
    var $ = layui.jquery;

    form.on('submit(formDemo)', function (data) {
        if (data.field.salt === '') {
            layer.msg('加密盐不能为空');
            return false
        }

        if (data.field.account === '') {
            layer.msg('管理员账号不能为空');
            return false
        }

        if (data.field.password === '') {
            layer.msg('密码不能为空');
            return false
        }

        if (data.field.password !== data.field.password_confirm) {
            layer.msg('密码不一致');
            return false
        }


        var d = {};
        d.BlogName = data.field.blog_name;
        d.AdminPageName = data.field.admin_page_name;
        d.Email = data.field.email;
        d.Salt = data.field.salt;
        d.Account = data.field.account;
        d.Password = data.field.password;

        $.ajax({
            type: "POST",
            contentType: "application/json",
            url: "/init",
            context: document.body,
            data: JSON.stringify(d),
            dataType: 'json',
            success: function (result) {
                console.log("result: " + result);
                layer.msg('初始化成功');
                setTimeout(function () {
                    window.location.href = "/admin";
                }, 1000);
            },
            error: function (result) {
                console.error(result);
                layer.msg('初始化失败');
            }
        });

        return false;
    });
});