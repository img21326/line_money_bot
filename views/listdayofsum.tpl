{{define "title"}}
Index
{{end}}

{{define "content"}}
<div class="text-center loading">
    <a>Loading</a>
</div>
<div class="mainbox">
    <ul class="nav justify-content-center">
        <li class="nav-item">
            <a class="nav-link date-desc" aria-current="page" href="#">
                << </a>
        </li>
        <li class="nav-item">
            <a class="nav-link" href="#" id="date"></a>
        </li>
        <li class="nav-item">
            <a class="nav-link date-insc" href="#">>></a>
        </li>
        <!-- <li class="nav-item">
        <a class="nav-link disabled" href="#" tabindex="-1" aria-disabled="true">Disabled</a>
    </li> -->
    </ul>
    <canvas id="myChart" width="300" height="300"></canvas>
    <ul class="list-group mt-3 listbox">

    </ul>
</div>


<script>
    var ctx = document.getElementById('myChart').getContext('2d');
    var dt = new Date();
    const color = [{
            bdc: 'rgba(255, 99, 132, 1)',
            bgc: 'rgba(255, 99, 132, 0.2)',
        },
        {
            bdc: 'rgba(54, 162, 235, 1)',
            bgc: 'rgba(54, 162, 235, 0.2)',
        },
        {
            bdc: 'rgba(255, 206, 86, 1)',
            bgc: 'rgba(255, 206, 86, 0.2)',
        },
        {
            bdc: 'rgba(75, 192, 192, 1)',
            bgc: 'rgba(75, 192, 192, 0.2)',
        },
        {
            bdc: 'rgba(153, 102, 255, 1)',
            bgc: 'rgba(153, 102, 255, 0.2)',
        },
    ]

    var config = {
        type: 'line',
        data: {
            labels: ['money'],
            datasets: []
        },
        options: {
            responsive: true,
            plugins: {
                legend: {
                    position: 'top',
                },
                title: {
                    display: true,
                    text: '單日總額'
                }
            }
        },
    };
    var myChart = new Chart(ctx, config);
    var userId = "";
    var cate = "{{.cate}}";
    // alert(cate);
    $(function() {
        $(".mainbox").hide();
        var liffID = '{{.liff_id}}';
        console.log(liffID);

        liff.init({
            liffId: liffID
        }).then(function() {
            console.log('LIFF init');
            liff.getProfile().then(user => {
                console.log(user.userId);
                userId = user.userId;
                // getData(user.userId, dt.getFullYear(), dt.getMonth() + 1)
                fetchAll(user.userId, dt.getFullYear(), dt.getMonth() + 1);
            })
        }).catch(function(error) {
            console.log(error);
        });
    });
    $(".date-desc").click(() => {
        dt = new Date(dt.setMonth(dt.getMonth() - 1));
        fetchAll(userId, dt.getFullYear(), dt.getMonth() + 1);
    });

    $(".date-insc").click(() => {
        dt = new Date(dt.setMonth(dt.getMonth() + 1));
        fetchAll(userId, dt.getFullYear(), dt.getMonth() + 1);
    });

    function loading() {
        $(".loading").fadeIn()
        $(".mainbox").fadeOut()
    }

    function closeloading() {
        $(".loading").fadeOut()
        $(".mainbox").fadeIn()
    }

    function fetchAll(userId, year, month) {
        $("#date").html(`${year}-${month} (${cate})`);
        loading();
        Promise.all([getData(userId, year, month)]).then(values => {
            $(".listbox").html("");
            config.data.datasets = [];
            $.each(values[0], (index, data) => {
                data.day = new Date(data.day);
            });
            var byDate = values[0].slice(0)
            byDate.sort(function(a, b) {
                return a.day - b.day;
            })
            var labels = [];
            var values = [];
            $.each(byDate, (index, data) => {
                console.log(data.day);
                $(".listbox").append(
                    `<li class="list-group-item d-flex justify-content-between align-items-center">
${data.day.getFullYear()}-${data.day.getMonth() + 1}-${data.day.getDate()}
<span class="badge bg-primary rounded-pill">$ ${data.total}</span>
                                                    </li>`
                );
                labels.push(`${data.day.getFullYear()}-${data.day.getMonth() + 1}-${data.day.getDate()}`);
                values.push(data.total);
            });
            console.log(values);
            config.data.labels = labels;
            config.data.datasets.push({
                label: 'Money',
                data: values,
                fill: true,
                borderColor: 'rgb(75, 192, 192)',
                tension: 0.1
            })
            myChart.update()
            closeloading()
        })
    }

    function getData(userId, year, month) {
        return new Promise(function(resolve, reject) {
            $.ajax({
                url: "/v1/acc/days/list/sum",
                type: "POST",
                cache: false,
                dataType: 'json',
                data: JSON.stringify({
                    "user_id": userId,
                    "year": year,
                    "month": month,
                    "cate": cate,
                }),
                contentType: "application/json",
                success: (res) => {
                    resolve(res);
                },

                error: function(xhr, ajaxOptions, thrownError) {
                    console.log(xhr.status);
                    console.log(thrownError);
                    reject(xhr);
                }
            });
        })
    }

    // document.getElementById("myChart").onclick = function(evt) {
    // var activePoints = myChart.getElementsAtEventForMode(evt, 'point', myChart.options);
    // var firstPoint = activePoints[0];
    // var clickCateName = config.data.datasets[firstPoint.datasetIndex].label;
    // var label = myChart.data.labels[firstPoint._index];
    // var value = myChart.data.datasets[firstPoint._datasetIndex].data[firstPoint._index];
    // alert(label + ": " + value);
    // };
</script>
{{end}}