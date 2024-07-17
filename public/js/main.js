function ActionListener() {
    var data = {
        row1: {
            a: '',
            b: '',
            c: '',
        },
        row2: {
            a: '',
            b: '',
            c: '',
        },
        row3: {
            a: '',
            b: '',
            c: '',
        },
    }

    cells = $('.column')
    var turn = 'X'

    cells.each(function(index) {
        $(this).on('click', function () {
            className = $(this).attr('class')
            row = className.slice(18,22);
            cell = className.slice(23);

            data[row][cell] = turn;

            if ($(this).text() != 'N') {
                return
            }

            $(this).html(turn)
            $(this).css('color','black')

            if (turn == 'X') {
                turn = 'O';
            } else {
                turn = 'X'
            }

            $('.turn').html(turn)
        })
    })
}

function WebSocketConn() {
    var conn;
    var btn = document.getElementById('websocket')

    btn.addEventListener('click', function() {
        conn.send('server')
    })

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws");
        conn.onclose = function(evt) {
            console.log('connection closed!');
        }
        conn.onmessage = function(evt) {
            console.log('have a response');
            console.log(evt.data);
        }
    } else {
        console.log('doest support websocket');
    }
}

document.addEventListener('DOMContentLoaded', function() {
    ActionListener()
    // WebSocketConn()
}, false);