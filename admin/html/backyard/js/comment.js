var commentTable;
layui.use(['table', 'layer'], function () {
    var table = layui.table;
    var layer = layui.layer;

    commentTable = table.render({
        elem: '#comment-list',
        url: '/admin/comments/list',
        page: false,
        cols: [[
            {field: 'ID', title: 'ID', width: 80, sort: true},
            {field: 'PostID', title: 'PostID'},
            {field: 'Content', title: 'Content'},
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

    table.on('tool(comment-list)', function (obj) {
        var data = obj.data;
        var layEvent = obj.event;

        switch (layEvent) {
            case 'edit':
                editComment(layer, obj);
                break;
            case 'open':
                window.open('/posts/' + obj.data.PostID, '_blank');
                break;
            case 'del':
                confirmDel(layer, data.ID, obj);
                break;
        }
    });
});

function editComment(layer, obj) {
    layer.prompt({
        formType: 2,
        value: obj.data.Content,
        title: '请输入评论内容',
    }, function (value, index, elem) {
        $.ajax({
            type: "PATCH",
            contentType: "application/json",
            url: "/admin/comments/" + obj.data.ID,
            context: document.body,
            data: JSON.stringify({
                "Content": value,
            }),
            dataType: 'json',
            success: function (result) {
                console.log("result: " + result);
                obj.update({
                    "Content": value
                });
                toast('更新评论成功');
            },
            error: function (result) {
                console.error(result);
                toast('更新评论失败');
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
        delComment(id, obj);
    });
}

function delComment(id, obj) {
    $.ajax({
        type: "DELETE",
        contentType: "application/json",
        url: "/admin/comments/" + id,
        context: document.body,
        success: function (result) {
            obj.del();
            toast('删除评论成功');
        },
        error: function (result) {
            console.error(result);
            toast('删除评论失败');
        }
    });
}