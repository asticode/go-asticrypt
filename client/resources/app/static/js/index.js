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
                case "email.added":
                    index.listenEmailAdded(message);
                    break;
                case "email.listed":
                    index.listenEmailListed(message);
                    break;
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
    listenEmailAdded: function(message) {
        asticode.modaler.hide();
        asticode.notifier.success(message.payload);
        index.sendEmailList();
    },
    listenEmailListed: function(message) {
        // Init content
        let content = `<div class="index-list">
            <div class="index-button">
                <button class="btn btn-success" onclick="index.onClickAddEmail()">Add a new email</button>
            </div>`;

        // Loop through emails
        for (let i = 0; i < message.payload.length; i++) {
            content += `<div class="index-item">
                ` + message.payload[i] + `
            </div>`;
        }

        // Close content
        content += "</div>";

        // Set content
        document.getElementById("index").innerHTML = content;
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
                index.sendEmailList();
                break;
            case "login":
                document.getElementById("index").innerHTML = `<div class="index-table">
                    <div class="index-cell">
                        <div class="index-form">
                            <input type="password" placeholder="Password" id="value-password" onkeypress="if (event.keyCode === 13) document.getElementById('btn').click()" autofocus>
                            <button class="btn btn-success btn-lg" id="btn" onclick="index.onClickLogin()">Login</button>
                        </div>
                    </div>
                </div>`;
                break;
            default:
                document.getElementById("index").innerHTML = `<div class="index-table">
                    <div class="index-cell">
                        <div class="index-form">
                            <input type="password" placeholder="Password" id="value-password" onkeypress="if (event.keyCode === 13) document.getElementById('btn').click()" autofocus>
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
    onClickAddEmail: function() {
        // Build content
        let content = document.createElement("div");
        content.innerHTML = `<input type="email" placeholder="Email" id="value-email" onkeypress="if (event.keyCode === 13) document.getElementById('btn').click()">
        <button class="btn btn-success btn-lg" id="btn" onclick="index.onClickSubmitEmail()">Add</button>`;

        // Update modal
        asticode.modaler.setContent(content);
        asticode.modaler.show();
        document.getElementById("value-email").focus();
    },
    onClickLogin: function() {
        index.sendIndexLogin(document.getElementById("value-password").value);
    },
    onClickSignup: function() {
        index.sendIndexSignup(document.getElementById("value-password").value);
    },
    onClickSubmitEmail: function() {
        index.sendEmailAdd(document.getElementById("value-email").value);
    },
    sendEmailAdd: function(email) {
        asticode.loader.show();
        astilectron.send({name: "email.add", payload: email});
    },
    sendEmailList: function() {
        asticode.loader.show();
        astilectron.send({name: "email.list"});
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