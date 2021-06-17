{{define "title"}}
Index
{{end}}

{{define "content"}}
<div class="text-center loading">
    <a>Loading</a>
</div>
<div class="mainbox">
    <select class="form-select" aria-label="Default select example">
    </select>
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
    var cate_arr = [];

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
                getListCate(userId).then(val => {
                    cate_arr = val;
                    $.each(val, (index, v) => {
                        if (cate != "" && cate == v) {
                            $('.form-select').append(`<option selected="selected" value="${v}">
                            ${v}<`+`/option>`); 
                        } else {
                            $('.form-select').append(`<option value="${v}">
                            ${v}<`+`/option>`);
                        }
                    })
                    if (cate == "") {
                        cate = cate_arr[0]
                    }
                    fetchAll(user.userId, dt.getFullYear(), dt.getMonth() + 1);
                })
                // getData(user.userId, dt.getFullYear(), dt.getMonth() + 1)
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
    $('.form-select').on('change', (event) => {
        cate = event.target.value
        console.log(cate)
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
        $("#date").html(`${year}-${month}`);
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
<span class="badge rounded-pill bg-primary">${data.day.getFullYear()}-${addZero(data.day.getMonth()+1)}-${addZero(data.day.getDate())}</span>
<span class="badge rounded-pill bg-secondary">$${data.total}</span>
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
                );
                labels.push(`${data.day.getFullYear()}-${addZero(data.day.getMonth() + 1)}-${addZero(data.day.getDate())}`);
                values.push(data.total);
            });
            config.data.labels = labels;

            config.data.datasets.push({
                label: '今日額度',
                data: values,
                fill: true,
                borderColor: 'rgba(255, 99, 132, 1)',
                backgroundColor: 'rgba(255, 99, 132, 0.2)',
                type: 'bar',
                stack: 'combined',
            })

            var sum = 0
            var nv = []
            $.each(values, (i, v) => {
                sum += v
                nv.push(sum)
            })

            config.data.datasets.push({
                label: '總額度',
                data: nv,
                borderColor: 'rgba(255, 99, 132, 1)',
                stack: 'combined',
            })
            myChart.update()
            closeloading()
        })
    }

    function addZero(s) {
        if (s.toString().length == 1) {
            return "0" + s.toString()
        }
        return s
    }

    function printTags(tags) {
        var html = ""
        $.each(tags, (index, tag) => {
            html += `<span class="badge rounded-pill bg-info text-dark">${tag}</span>`
        })
        return html
    }

    $("#accordionPanelsStayOpenExample").on("click", ".accordion-item", function(event) {
        if ($(this).find(".accordion-body").attr("data-status") == "-1") {
            var year = parseInt($(this).attr("data-year"));
            var month = parseInt($(this).attr("data-month"));
            var day = parseInt($(this).attr("data-date"));
            getListDayOfInfoData(userId, year, month, day).then(value => {
                $(this).find(".accordion-body").attr("data-status", "1")
                $(this).find(".accordion-body").html("")
                if (value.length > 0) {
                    html_ =
                        `<table class="table" style="width:100%;word-break:break-all; word-wrap:break-all;"><thead><tr><th scope="col">Date</th><th scope="col">Total</th><th scope="col">Tags</th></tr></thead><tbody>`
                    $.each(value, (index, data) => {
                        var d = new Date(data.created_at)
                        html_ += `<tr>
<th scope="row" style="width:30%;">${addZero(d.getHours())}:${addZero(d.getMinutes())}</th>
<td style="width:30%;">${data.amount}</td>
<td style="width:40%;">${printTags(data.tags)}</td>
                             </tr>`

                    })
                    $(this).find(".accordion-body").html(html_)
                }
            })
        }

    });

    function getListCate(userId) {
        return new Promise(function(resolve, reject) {
            $.ajax({
                url: "/v1/acc/month/list/cate",
                type: "POST",
                cache: false,
                dataType: 'json',
                data: JSON.stringify({
                    "user_id": userId,
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