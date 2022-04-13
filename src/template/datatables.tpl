
<!DOCTYPE html>
<html>
<head>
<link rel="shortcut icon" href="static/favicon.ico" type="image/x-icon">
<link rel="icon" href="static/favicon.ico" type="image/x-icon">
<title>{{ .Title }}</title>
</head>
<body>
    <link rel="stylesheet" href="https://cdn.datatables.net/1.10.19/css/dataTables.bootstrap4.min.css">
    <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/twitter-bootstrap/4.1.1/css/bootstrap.css">

    <script src="https://code.jquery.com/jquery-3.3.1.js"></script>
    <script src="https://cdn.datatables.net/1.10.19/js/jquery.dataTables.min.js"></script>
    <script src="https://cdn.datatables.net/1.10.19/js/dataTables.bootstrap4.min.js"></script>
    <script src="https://cdn.jsdelivr.net/gh/jeffreydwalter/ColReorderWithResize@9ce30c640e394282c9e0df5787d54e5887bc8ecc/ColReorderWithResize.js"></script>


<br>
<h1 align='center'>{{ .Title }}</h1>
<br>

<div class="container">
        <table id="example"class="table table-striped table-bordered" style="width:100%">
            <thead class="thead-dark">
                <tr>            
                    {{ range .Columns }}
                        
                    <th scope="col">{{.}}</th>
                    {{ end }}
                </tr>
            </thead>
        </table>
    </div>

<script>
    $(document).ready(function() {
        $('#example').DataTable( {
        "responsive": true,
        'dom': 'Rlfrtip',
        "aaData":[{{range $i,$row := .Rows}}{{if $i}},{{end}}{ {{ range $index,$element := $row }}{{if $index}},{{end}}"{{index $.Columns $index}}":"{{ replace $element "\n" "<br />" }}"{{ end }}} {{end}}],
        "aoColumns":
        [
            {{ range .Columns }}
            {"mDataProp": "{{.}}", className: "text-center"},
            {{ end }}
        ]
        } );
    } );

    </script>
</body>
</html>