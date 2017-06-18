var index = {
    init: function() {
        // Init
        asticode.loader.init();
        asticode.notifier.init();
        asticode.modaler.init();

        // Wait for astilectron to be ready
        document.addEventListener('astilectron-ready', function() {
            // Listen
            index.listen();

            // Send index.show
            index.sendIndexShow();
        })
    },
    listen: function() {
        astilectron.listen(function(message) {
            asticode.loader.hide();
            switch (message.name) {
                case "error":
                    index.listenError(message);
                    break;
                case "index.logged.in":
                    index.listenIndexLoggedIn();
                    break;
                case "index.show":
                    index.listenIndexShow(message);
                    break;
                case "index.signed.up":
                    index.listenIndexSignedUp();
                    break;
            }
        });
    },
    listenError: function(message) {
        asticode.notifier.error(message.payload);
    },
    listenIndexLoggedIn: function() {
        index.sendIndexShow();
    },
    listenIndexShow: function(message) {
        switch (message.payload) {
            case "index":
                document.getElementById("index").innerHTML = `Index`;
                break;
            case "login":
                document.getElementById("index").innerHTML = `<div class="index-table">
                    <div class="index-cell">
                        <div class="index-form">
                            <input type="password" placeholder="Password" id="value-password" onkeypress="if (event.keyCode == 13) document.getElementById('btn').click()" autofocus>
                            <button class="btn btn-success btn-lg" id="btn" onclick="index.onClickLogin()">Login</button>
                        </div>
                    </div>
                </div>`;
                break;
            default: // Sign up
                document.getElementById("index").innerHTML = `<div class="index-table">
                    <div class="index-cell">
                        <div class="index-form">
                            <input type="password" placeholder="Password" id="value-password" onkeypress="if (event.keyCode == 13) document.getElementById('btn').click()" autofocus>
                            <button class="btn btn-success btn-lg" id="btn" onclick="index.onClickSignup()">Sign up</button>
                        </div>
                    </div>
                </div>`;
                break;
        }
    },
    listenIndexSignedUp: function() {
        index.sendIndexShow();
    },
    onClickLogin: function() {
        index.sendIndexLogin(document.getElementById("value-password").value);
    },
    onClickSignup: function() {
        index.sendIndexSignup(document.getElementById("value-password").value);
    },
    sendIndexLogin: function(password) {
        asticode.loader.show();
        astilectron.send({name: "index.login", payload: password});
    },
    sendIndexShow: function() {
        asticode.loader.show();
        astilectron.send({name: "index.show"});
    },
    sendIndexSignup: function(password) {
        asticode.loader.show();
        astilectron.send({name: "index.sign.up", payload: password});
    }
};