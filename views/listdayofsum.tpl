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
    <div class="accordion" id="accordionPanelsStayOpenExample">

    </div>
</div>


<script>
    var ctx = document.getElementById('myChart').getContext('2d');
    var dt = new Date();

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
                    text: '每日總額'
                }
            }
        },
    };
    var myChart = new Chart(ctx, config);
    var userId = "";
    var cate = "{{.cate}}";

    $(function() {
        $(".mainbox").hide();
        var liffID = '{{.liff_id}}';
        console.log(liffID);
        // fetchAll(userId, dt.getFullYear(), dt.getMonth() + 1);


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
        Promise.all([getListDayOfSumData(userId, year, month)]).then(values => {
            $("#accordionPanelsStayOpenExample").html("");
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
                $("#accordionPanelsStayOpenExample").append(
                    `<div class="accordion-item" data-year="${data.day.getFullYear()}" data-month="${data.day.getMonth()+1}" data-date="${data.day.getDate()}"  data-cate="${cate}">` +
                    `<h2 class="accordion-header" id="panelsStayOpen-heading${index}">` +
                    `<button class="accordion-button collapsed" type="button" data-bs-toggle="collapse"
data-bs-target="#panelsStayOpen-collapse${index}" aria-expanded="false"
aria-controls="panelsStayOpen-collapse${index}">
<span class="badge rounded-pill bg-info text-dark">${data.day.getFullYear()}-${data.day.getMonth()+1}-${data.day.getDate()}</span>
<span class="badge rounded-pill bg-light text-dark">$${data.total}</span>
                                    </button>
                                </h2>
<div id="panelsStayOpen-collapse${index}" class="accordion-collapse collapse"
aria-labelledby="panelsStayOpen-heading${index}">
                                    <div class="accordion-body" data-status="-1">
                                    <div class="d-flex justify-content-center">
                    <div class="spinner-border" role="status">
                        <span class="visually-hidden">Loading...</span>
                    </div>
                    </div> 
                                    </div>
                                </div>
                            </div>`

                    // <table class="table">
                    //     <thead>
                    //         <tr>
                    //             <th scope="col">Date</th>
                    //             <th scope="col">Total</th>
                    //             <th scope="col">Tags</th>
                    //         </tr>
                    //     </thead>
                    //     <tbody>
                    //         <tr>
                    //             <th scope="row">1</th>
                    //             <td>Mark</td>
                    //             <td>Otto</td>
                    //         </tr>
                    //         <tr>
                    //             <th scope="row">2</th>
                    //             <td>Jacob</td>
                    //             <td>Thornton</td>
                    //         </tr>
                    //     </tbody>
                    // </table>
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
                borderColor: 'rgba(255, 99, 132, 1)',
                backgroundColor: 'rgba(255, 99, 132, 0.2)',
                tension: 0.1
            })
            myChart.update()
            closeloading()
        })
    }

    $("#accordionPanelsStayOpenExample").on("click", ".accordion-item", function(event) {
        if ($(this).find(".accordion-body").attr("data-status") == "-1") {
            var year = parseInt($(this).attr("data-year"));
            var month = parseInt($(this).attr("data-month"));
            var day = parseInt($(this).attr("data-date"));
            getListDayOfInfoData(userId, year, month, day).then(value => {
                $(this).find(".accordion-body").attr("data-status", "1")
                $(this).find(".accordion-body").html("")
                console.log(value.length)
                console.log(value.length > 0)
                if (value.length > 0) {
                    html_ =
                        `<table class="table"><thead><tr><th scope="col">Date</th><th scope="col">Total</th><th scope="col">Tags</th></tr></thead><tbody>`
                    $.each(value, (index, data) => {
                        var d = new Date(data.created_at)
                        html_ += `<tr>
<th scope="row">${d.getHours()}:${d.getMinutes()}</th>
<td>${data.amount}</td>
<td>${data.tags.join()}</td>
                             </tr>`

                    })
                    console.log(html_)
                    $(this).find(".accordion-body").html(html_)
                }
            })
        }

    });

    function getListDayOfInfoData(userId, year, month, day) {
        return new Promise(function(resolve, reject) {
            $.ajax({
                url: "/v1/acc/day/list/info",
                type: "POST",
                cache: false,
                dataType: 'json',
                data: JSON.stringify({
                    "user_id": userId,
                    "year": year,
                    "month": month,
                    "day": day,
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

    function getListDayOfSumData(userId, year, month) {
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