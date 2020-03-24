layui.use(['table', 'layer'], function () {
    var table = layui.table;
    var layer = layui.layer;

    table.render({
        elem: '#post-list',
        url: '/admin/posts/list',
        page: true,
        cols: [[
            {field: 'ID', title: 'ID', width: 80, sort: true, fixed: 'left'},
            {field: 'Title', title: 'Title'},
            {field: 'Description', title: 'Description'},
            {field: 'Cover', title: 'Cover'},
            {field: 'Content', title: 'Content'},
            {field: 'Created', title: 'Created', sort: true, width: 160, templet: formatCreated},
            {field: 'Updated', title: 'Updated', sort: true, width: 160, templet: formatUpdated},
            {field: 'IsPublished', title: 'IsPublished', width: 80, sort: true},
            {field: 'Type', title: 'Type', width: 100, sort: true, templet: formatType},
            {field: 'Tags', title: 'Tags', width: 150, templet: addTags},
            {fixed: 'right', align: 'left', width: 170, toolbar: '#actions'}
        ]],
        parseData: function (res) {
            return {
                "code": res.ok ? 0 : -1,
                "msg": res.msg,
                "count": res.total,
                "data": res.result
            };
        }
    });

    table.on('tool(post-list)', function (obj) {
        var data = obj.data;
        var layEvent = obj.event;

        switch (layEvent) {
            case 'detail':
                window.open('/posts/' + data.ID, '_blank');
                break;
            case 'edit':
                window.location.href = '/admin/posts/editor/' + data.ID;
                break;
            case 'del':
                confirmDel(layer, data.ID, obj);
                break;
        }
    })
});

function confirmDel(layer, id, obj) {
    //具体配置参考：http://www.layui.com/doc/modules/layer.html
    layer.confirm('你确定要删除这条数据吗？', {
        btn: ['按错了', '删除']
    }, function (index, layero) {
        layer.closeAll();
    }, function (index) {
        layer.closeAll();
        delPost(id, obj);
    });
}

function delPost(id, obj) {
    $.ajax({
        type: "DELETE",
        contentType: "application/json",
        url: "/admin/posts/" + id,
        context: document.body,
        success: function (result) {
            obj.del();
            toast('删除文章成功');
        },
        error: function (result) {
            console.error(result);
            toast('删除文章失败');
        }
    });
}

function addTags(d) {
    var r = "";
    d.Tags.forEach(function (item, index) {
        r += `<span class="layui-badge layui-bg-blue">${item.Name}</span>\n`
    });
    return r
}

function formatType(d) {
    switch (d.Type) {
        case 1:
            return "article";
        case 2:
            return "snippet";
        case 3:
            return "moment";
        default:
            return "unknown";
    }
}