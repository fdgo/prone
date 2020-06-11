<!DOCTYPE html>
<html>

<head>
    {{template "header" .node}}
    <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/morris.js/0.5.1/morris.css">
    <link rel="stylesheet" href="//cdnjs.cloudflare.com/ajax/libs/bootstrap-fileinput/4.5.2/css/fileinput.min.css"/>
</head>

<body class="hold-transition skin-blue sidebar-mini">
<div class="wrapper">
    {{template "navbar" .}}
    <div class="content-wrapper">
        {{template "title" .}}
        <section class="content">
            <div class="row">
                <div class="col-lg-3 col-xs-6">
                    <div class="small-box bg-aqua">
                        <div class="inner">
                            <h3>150</h3>
                            <p>新订单</p>
                        </div>
                        <div class="icon">
                            <i class="ion ion-bag"></i>
                        </div>
                        <a href="#" class="small-box-footer" data-target="#upload-modal" data-toggle="modal">更多 <i
                                    class="fa fa-arrow-circle-right"></i></a>
                    </div>
                </div>
                <div class="col-lg-3 col-xs-6">
                    <div class="small-box bg-green">
                        <div class="inner">
                            <h3>53<sup style="font-size: 20px">%</sup></h3>
                            <p>增长率</p>
                        </div>
                        <div class="icon">
                            <i class="ion ion-stats-bars"></i>
                        </div>
                        <a href="#" class="small-box-footer">更多 <i class="fa fa-arrow-circle-right"></i></a>
                    </div>
                </div>
                <div class="col-lg-3 col-xs-6">
                    <div class="small-box bg-yellow">
                        <div class="inner">
                            <h3>44</h3>
                            <p>用户注册</p>
                        </div>
                        <div class="icon">
                            <i class="ion ion-person-add"></i>
                        </div>
                        <a href="#" class="small-box-footer">更多 <i class="fa fa-arrow-circle-right"></i></a>
                    </div>
                </div>
                <div class="col-lg-3 col-xs-6">
                    <div class="small-box bg-red">
                        <div class="inner">
                            <h3>65</h3>
                            <p>访问量</p>
                        </div>
                        <div class="icon">
                            <i class="ion ion-pie-graph"></i>
                        </div>
                        <a href="#" class="small-box-footer">更多 <i class="fa fa-arrow-circle-right"></i></a>
                    </div>
                </div>
            </div>
            <div class="row">
                <section class="col-lg-7">
                    <div class="box box-warning" id="area-chart" data-source="?action=user">
                        <div class="box-header with-border">
                            <h3 class="box-title">趋势</h3>
                            <div class="box-tools pull-right">
                                <a class="btn btn-box-tool" data-widget="collapse"><i class="fa fa-minus"></i></a>
                                <a class="btn btn-box-tool refresh-btn"><i class="fa fa-refresh"></i></a>
                            </div>
                        </div>
                        <div class="box-body">
                            <div class="chart" style="height: 250px;"></div>
                        </div>
                    </div>
                </section>
            </div>
        </section>
    </div>

    <div class="modal" id="upload-modal">
        <div class="modal-dialog">
            <div class="modal-content box">
                <form action="/" method="post" class="form-horizontal">
                    <div class="modal-header">
                        <a class="close" data-dismiss="modal" aria-label="Close"><span aria-hidden="true">&times;</span></a>
                        <h4 class="modal-title">上传文件</h4>
                    </div>
                    <div class="modal-body">
                        <div class="form-group">
                            <label class="col-sm-3 control-label">MD5</label>
                            <div class="col-sm-6">
                                <input type="text" class="form-control" name="md5" readonly>
                            </div>
                        </div>
                        <div class="form-group">
                            <label class="col-sm-3 control-label">路径</label>
                            <div class="col-sm-7">
                                <input type="text" class="form-control" name="url" readonly>
                            </div>
                        </div>
                        <div class="form-group">
                            <label class="col-sm-3 control-label">上传</label>
                            <div class="col-sm-8">
                                <input type="file" class="file" data-upload-url="/upload" data-show-preview="false"
                                       data-show-remove="false"
                                       data-language="zh">
                            </div>
                        </div>
                    </div>
                    <div class="modal-footer">
                        <a class="btn btn-default" data-dismiss="modal">取消</a>
                        <button type="submit" class="btn bg-purple">保存</button>
                    </div>
                </form>
            </div>
        </div>
    </div>
    {{template "footer"}}
    <script src="//cdnjs.cloudflare.com/ajax/libs/raphael/2.1.0/raphael-min.js"></script>
    <script src="//cdnjs.cloudflare.com/ajax/libs/morris.js/0.5.1/morris.min.js"></script>
    <script type="text/javascript">
        $(document).ready(function () {
            var area = new Morris.Line({
                element: $('#area-chart .chart'),
                xkey: 'y',
                hideHover: 'auto',
                ykeys: ['item1', 'item2'],
                labels: ['类型1', '类型2'],
                lineColors: ['#a0d0e0', '#3c8dbc']
            });
            $('#area-chart').boxRefresh({
                loadInContent: false,
                responseType: 'json',
                onLoadDone: function (e) {
                    area.setData(e.data);
                }
            })
        })
    </script>

    <script src="//cdnjs.cloudflare.com/ajax/libs/bootstrap-fileinput/4.5.2/js/fileinput.min.js"></script>
    <script src="//cdnjs.cloudflare.com/ajax/libs/bootstrap-fileinput/4.5.2/js/locales/zh.min.js"></script>
    <script type="text/javascript">
        $(document).on('fileuploaded', function (ev, d) {
            var resp = d.response;
            if (resp.code == 200)
                $('.modal :text[readonly]').each(function (i, el) {
                    el.value = resp.data[el.name];
                })
        })
    </script>
</div>
</body>

</html>