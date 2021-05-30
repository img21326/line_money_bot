{{define "title"}}
Index
{{end}}

{{define "content"}}
<ul class="nav justify-content-center">
    <li class="nav-item">
        <a class="nav-link active" aria-current="page" href="#">
            << </a>
    </li>
    <li class="nav-item">
        <a class="nav-link" href="#" id="date"></a>
    </li>
    <li class="nav-item">
        <a class="nav-link" href="#">>></a>
    </li>
    <!-- <li class="nav-item">
        <a class="nav-link disabled" href="#" tabindex="-1" aria-disabled="true">Disabled</a>
    </li> -->
</ul>
<canvas id="myChart" width="300" height="300"></canvas>
<ul class="list-group mt-3 listbox">

</ul>

<script>
    var ctx = document.getElementById('myChart').getContext('2d');
    var dt = new Date();
    $("#date").html(dt.getFullYear() + "-" + (dt.getMonth() +1))
    const color = [
        {
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
            datasets: [
                // {
                //     label: 'Dataset 1',
                //     data: [1],
                //     borderColor: 'rgba(255, 99, 132, 1)',
                //     backgroundColor: 'rgba(255, 99, 132, 0.2)',
                // },
                // {
                //     label: 'Dataset 2',
                //     data: [-20],
                //     borderColor: 'rgba(54, 162, 235, 1)',
                //     backgroundColor: 'rgba(54, 162, 235, 0.2)',
                // }
            ]
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
    $(function () {
        var liffID = '1656043897-1kg8a3DM';

        liff.init({
            liffId: liffID
        }).then(function () {
            console.log('LIFF init');
            liff.getProfile().then(user => {
                console.log(user.userId)
                getData(user.userId, dt.getFullYear(), dt.getMonth() +1)
            })
        }).catch(function (error) {
            console.log(error);
        });
    });

    function getData(userId, year, month) {
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
                        $.each(res, (index, data) => {
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
                    },

                    error: function (xhr, ajaxOptions, thrownError) {
                        console.log(xhr.status);
                        console.log(thrownError);
                    }
                });
    }
</script>
{{end}}