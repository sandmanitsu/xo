var data = {
    turn: '',
    userTurn: '',
    username: '',
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

// ws connect
var conn;

var turn;
var roomId;
var clientId;
var username;

function ActionListener() {
    cells = $('.column')

    cells.each(function(index) {
        $(this).on('click', function () {
            // console.log(data);
            // if (data.userTurn != username) {
            //     console.log('locked for ' + username);
            //     return;
            // }

            className = $(this).attr('class')
            row = className.slice(7,11);
            cell = className.slice(12);

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

            data.turn = turn
            data.username = username

            conn.send(JSON.stringify(data))

            $('.turn').html(turn)
        })
    })
}

function startGame(message) {
    if (message == 'A new user join the room') {
        $.ajax({
            method: 'get',
            url: '/ws/getClients/' + roomId,
            success: function(clients) {
                clientsList = clients

                if (clients.length == 2) {
                    $('.message-app').append('<div>Game started!' + getDateForMessage() + '</div>')
                    $('.message-app').append('<div>' + clients[0].user + ' you turn first! -_-' + getDateForMessage() + '</div>')

                    data.userTurn = clients[0].user

                    conn.send(JSON.stringify(data))
                }
            }
        });
    }
}

function syncPlayground() {
    if (data.row1.a) {
        $('.row1-a').html(data.row1.a)
        $('.row1-a').css('color','black')
    }
    if (data.row1.b) {
        $('.row1-b').html(data.row1.b)
        $('.row1-b').css('color','black')
    }
    if (data.row1.c) {
        $('.row1-c').html(data.row1.c)
        $('.row1-c').css('color','black')
    }

    if (data.row2.a) {
        $('.row2-a').html(data.row2.a)
        $('.row2-a').css('color','black')
    }
    if (data.row2.b) {
        $('.row2-b').html(data.row2.b)
        $('.row2-b').css('color','black')
    }
    if (data.row2.c) {
        $('.row2-c').html(data.row2.c)
        $('.row2-c').css('color','black')
    }

    if (data.row3.a) {
        $('.row3-a').html(data.row3.a)
        $('.row3-a').css('color','black')
    }
    if (data.row3.b) {
        $('.row3-b').html(data.row3.b)
        $('.row3-b').css('color','black')
    }
    if (data.row3.c) {
        $('.row3-c').html(data.row3.c)
        $('.row3-c').css('color','black')
    }
}

function getDateForMessage() {
    date = new Date();

    return ' ' + date.getHours() + ':' + date.getMinutes() + ':' + date.getSeconds()
}

function addNewUserHtml(username) {
    date = new Date();

    html = '<div>' + username + ' joined the room!' + getDateForMessage() + '</div>'

    $('.message-app').append(html)
}

function WebSocketConn() {
    roomId = window.location.pathname.split('/')[2]
    clientId = window.location.pathname.split('/')[3]
    username = window.location.pathname.split('/')[4]

    $('.user').append('<div>' + username + '</div>')

    if (window["WebSocket"]) {
        conn = new WebSocket("ws://" + document.location.host + "/ws/joinRoom/" + roomId + "?userId= " + clientId + "&user=" + username);
        conn.onclose = function(evt) {
            console.log('connection closed!');
        }
        conn.onmessage = function(evt) {
            console.log('have a response');

            response =  JSON.parse(evt.data)
            console.log(response);

            if (response.content == 'A new user join the room') {
                addNewUserHtml(response.user)
            }

            data = JSON.parse(response.playground)
            turn = data.turn
            $('.turn').html(turn)

            startGame(response.content);
            syncPlayground()

            console.log(data);
        }
    } else {
        console.log('doest support websocket');
    }
}

document.addEventListener('DOMContentLoaded', function() {
    ActionListener()
    WebSocketConn()
}, false);