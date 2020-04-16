layui.use('form', function () {
    var form = layui.form;

    $.ajax({
        type: "GET",
        contentType: "application/json",
        url: "/admin/preferences/get",
        context: document.body,
        success: function (result) {
            result = result.result;
            form.val("prefForm", {
                "blog_name": result.blog_name,
                "admin_page_name": result.admin_page_name,
                "about_id": result.about_id,
                "email": result.email,
                "qq_app_id": result.qq_app_id,
                "qq_app_key": result.qq_app_key,
                "qq_redirect": result.qq_redirect,
                "ta_id": result.ta_id,
                "beian": result.beian,
            });
            form.render();
        },
        error: function (result) {
            console.log(result);
            toast('加载列表失败');
        }
    });

    form.on('submit(prefForm)', function (data) {
        data = form.val('prefForm');
        $.ajax({
            type: "PATCH",
            contentType: "application/json",
            url: "/admin/preferences",
            context: document.body,
            data: JSON.stringify({
                "blog_name": data.blog_name,
                "admin_page_name": data.admin_page_name,
                "about_id": parseInt(data.about_id),
                "email": data.email,
                "qq_app_id": data.qq_app_id,
                "qq_app_key": data.qq_app_key,
                "qq_redirect": data.qq_redirect,
                "ta_id": data.ta_id,
                "beian": data.beian,
            }),
            success: function (result) {
                console.log(result);
                toast('修改配置成功');
            },
            error: function (result) {
                console.log(result);
                toast('修改配置成功');
            }
        });
        return false;
    })
});