var tagTable;
layui.use(['table', 'layer'], function () {
    var table = layui.table;
    var layer = layui.layer;

    tagTable = table.render({
        elem: '#tag-list',
        url: '/admin/tags/list',
        page: false,
        cols: [[
            {field: 'ID', title: 'ID', width: 80, sort: true},
            {field: 'Name', title: 'Name'},
            {field: 'Created', title: 'Created', sort: true, width: 160, templet: formatCreated},
            {field: 'Updated', title: 'Updated', sort: true, width: 160, templet: formatUpdated},
            {align: 'left', width: 170, toolbar: '#actions'}
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

    table.on('tool(tag-list)', function (obj) {
        var data = obj.data;
        var layEvent = obj.event;

        switch (layEvent) {
            case 'edit':
                editTag(layer, obj);
                break;
            case 'del':
                confirmDel(layer, data.ID, obj);
                break;
        }
    })

    $('#btn-create').click(function(){
        layer.prompt({
            formType: 2,
            title: '请输入标签名称',
        }, function (value, index, elem) {
            $.ajax({
                type: "POST",
                contentType: "application/json",
                url: "/admin/tags",
                context: document.body,
                data: JSON.stringify({
                    "Name": value,
                }),
                success: function (result) {
                    tagTable.reload();
                    toast('创建标签成功');
                },
                error: function (result) {
                    console.error(result);
                    toast('创建标签失败');
                }
            });
            layer.close(index);
        });
    });
});

function editTag(layer, obj) {
    layer.prompt({
        formType: 2,
        value: obj.data.Name,
        title: '请输入标签名称',
    }, function (value, index, elem) {
        $.ajax({
            type: "PATCH",
            contentType: "application/json",
            url: "/admin/tags/" + obj.data.ID,
            context: document.body,
            data: JSON.stringify({
                "Name": value,
            }),
            dataType: 'json',
            success: function (result) {
                console.log("result: " + result);
                obj.update({
                    "Name": value
                });
                toast('更新标签成功');
            },
            error: function (result) {
                console.error(result);
                toast('更新标签失败');
            }
        });
        layer.close(index);
    });
}

function confirmDel(layer, id, obj) {
    layer.confirm('你确定要删除这条数据吗？', {
        btn: ['按错了', '删除']
    }, function (index, layero) {
        layer.closeAll();
    }, function (index) {
        layer.closeAll();
        delTag(id, obj);
    });
}

function delTag(id, obj) {
    $.ajax({
        type: "DELETE",
        contentType: "application/json",
        url: "/admin/tags/" + id,
        context: document.body,
        success: function (result) {
            obj.del();
            toast('删除标签成功');
        },
        error: function (result) {
            console.error(result);
            toast('删除标签失败');
        }
    });
}