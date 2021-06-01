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
        type: 'bar',
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
                    text: '月總和'
                }
            }
        },
    };
    var myChart = new Chart(ctx, config);
    var userId = "";
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
        $("#date").html(year + "-" + month);
        loading()
        Promise.all([getData(userId, year, month),
            getTotal(userId, year, month)
        ]).then(values => {
            $(".listbox").html("")
            config.data.datasets = []
            $.each(values[0], (index, data) => {
                $(".listbox").append(
                    `<li class="list-group-item d-flex justify-content-between align-items-center">
${data.name}
<span class="badge bg-primary rounded-pill">$ ${data.total}</span>
                                                    </li>`
                );
                config.data.datasets.push({
                    label: data.name,
                    data: [data.total],
                    borderColor: color[index % 5].bdc,
                    backgroundColor: color[index % 5].bgc,
                })
                myChart.update()
            });
            $(".listbox").append(
                `<li class="list-group-item d-flex justify-content-between align-items-center list-group-item-secondary">
Total
<span class="badge bg-primary rounded-pill">$ ${values[1].total}</span>
                                                    </li>`
            );
            closeloading()
        })
    }

    function getData(userId, year, month) {
        return new Promise(function(resolve, reject) {
            $.ajax({
                url: "/v1/tags/sum",
                type: "POST",
                cache: false,
                dataType: 'json',
                data: JSON.stringify({
                    "user_id": userId,
                    "year": year,
                    "month": month,
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

    function getTotal(userId, year, month) {
        return new Promise(function(resolve, reject) {
            $.ajax({
                url: "/v1/user/total",
                type: "POST",
                cache: false,
                dataType: 'json',
                data: JSON.stringify({
                    "user_id": userId,
                    "year": year,
                    "month": month,
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
</script>
{{end}}