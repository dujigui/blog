layui.use('form', function () {
    var form = layui.form;

    $.ajax({
        type: "GET",
        contentType: "application/json",
        url: "/admin/preferences/get",
        context: document.body,
        success: function (result) {
            console.log(result);
            form.val("prefForm", {
                "blog_name": result.result.blog_name,
                "admin_page_name": result.result.admin_page_name,
                "about_id": result.result.about_id,
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