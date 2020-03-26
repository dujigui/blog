var fileTable;
layui.use(['upload','table', 'layer'], function () {
    var upload = layui.upload;
    var table = layui.table;
    var layer = layui.layer;

    upload.render({
        elem: '#btn-upload',
        url: '/admin/files',
        accept: 'file',
        done: function (res) {
            toast('上传文件成功');
            fileTable.reload();
        },
        error: function () {
            toast('上传文件失败');
        }
    });

    fileTable = table.render({
        elem: '#file-list',
        url: '/admin/files/list',
        page: true,
        cols: [[
            {field: 'ID', title: 'ID', width: 80, sort: true},
            {field: 'Original', title: 'Original'},
            {field: 'Hashed', title: 'Hashed'},
            {field: 'Created', title: 'Created', sort: true, width: 160, templet: formatCreated},
            {field: 'Updated', title: 'Updated', sort: true, width: 160, templet: formatUpdated},
            {title: 'Preview', width: 150, templet: preview},
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

    table.on('tool(file-list)', function (obj) {
        var data = obj.data;
        var layEvent = obj.event;

        switch (layEvent) {
            case 'open':
                window.open('/files/' + data.Hashed, '_blank');
                break;
            case 'url':
                var text = '/files/' + data.Hashed;
                clipboard(text);
                break;
            case 'del':
                confirmDel(layer, data.ID, obj);
                break;
        }
    })
});

function isImg(fn) {
    fn = fn.toLowerCase();
    return fn.endsWith(".jpg") || fn.endsWith(".png") || fn.endsWith(".jpeg") || fn.endsWith(".gif");
}

function preview(d) {
    if (isImg(d.Hashed)) {
        return `<img class="list-img" src="/files/${d.Hashed}"/>`;
    } else {
        return `<span>-</span>`;
    }
}

function clipboard(text) {
    var textarea = document.createElement("textarea");
    textarea.textContent = text;
    textarea.style.position = "fixed";
    document.body.appendChild(textarea);
    textarea.select();
    try {
        document.execCommand("copy");
    } catch (ex) {
        console.warn("Copy to clipboard failed.", ex);
    } finally {
        document.body.removeChild(textarea);
    }
}

function confirmDel(layer, id, obj) {
    layer.confirm('你确定要删除这条数据吗？', {
        btn: ['按错了', '删除']
    }, function (index, layero) {
        layer.closeAll();
    }, function (index) {
        layer.closeAll();
        delFile(id, obj);
    });
}

function delFile(id, obj) {
    $.ajax({
        type: "DELETE",
        contentType: "application/json",
        url: "/admin/files/" + id,
        context: document.body,
        success: function (result) {
            obj.del();
            toast('删除文件成功');
        },
        error: function (result) {
            console.error(result);
            toast('删除文件失败');
        }
    });
}