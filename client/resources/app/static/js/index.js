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

            // Send index
            index.sendIndex();
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
                case "indexed":
                    index.listenIndexed(message);
                    break;
                case "logged.in":
                    index.listenLoggedIn();
                    break;
                case "logged.out":
                    index.listenLoggedOut();
                    break;
                case "signed.up":
                    index.listenSignedUp();
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
        let content = `<div class="index-header">
            <button class="btn btn-success" onclick="index.onClickEmailAdd()" title="Add a new email"><i class="fa fa-plus"></i></button>
            <button class="btn btn-success" onclick="index.onClickEmailList()" title="Refresh emails list"><i class="fa fa-refresh"></i></button>
            <button class="btn btn-success" onclick="index.onClickLogout()" title="Log out"><i class="fa fa-sign-out"></i></button>
        </div>`;

        // Loop through emails
        content += `<div class="index-list">`;
        for (let i = 0; i < message.payload.length; i++) {
            content += `<div class="index-item">
                ` + message.payload[i] + `
            </div>`;
        }
        content += "</div>";

        // Set content
        document.getElementById("index").innerHTML = content;
    },
    listenError: function(message) {
        asticode.notifier.error(message.payload);
    },
    listenIndexed: function(message) {
        switch (message.payload) {
            case "index":
                index.sendEmailList();
                break;
            case "login":
                document.getElementById("index").innerHTML = `<div class="index-table">
                    <div class="index-cell">
                        <div class="index-form">
                            <input type="password" placeholder="Password" id="value-password" onkeypress="if (event.keyCode === 13) document.getElementById('btn-login').click()">
                            <button class="btn btn-success btn-lg" id="btn-login" onclick="index.onClickLogin()">Login</button>
                        </div>
                    </div>
                </div>`;
                document.getElementById("value-password").focus();
                break;
            default:
                document.getElementById("index").innerHTML = `<div class="index-table">
                    <div class="index-cell">
                        <div class="index-form">
                            <input type="password" placeholder="Password" id="value-password" onkeypress="if (event.keyCode === 13) document.getElementById('btn-signup').click()">
                            <button class="btn btn-success btn-lg" id="btn-signup" onclick="index.onClickSignUp()">Sign up</button>
                        </div>
                    </div>
                </div>`;
                document.getElementById("value-password").focus();
                break;
        }
    },
    listenLoggedIn: function() {
        index.sendIndex();
    },
    listenLoggedOut: function() {
        index.sendIndex();
    },
    listenSignedUp: function() {
        index.sendIndex();
    },
    onClickEmailAdd: function() {
        // Build content
        let content = document.createElement("div");
        content.innerHTML = `<input type="email" placeholder="Email" id="value-email" onkeypress="if (event.keyCode === 13) document.getElementById('btn-email').click()">
        <button class="btn btn-success btn-lg" id="btn-email" onclick="index.onClickSubmitEmail()">Add</button>`;

        // Update modal
        asticode.modaler.setContent(content);
        asticode.modaler.show();
        document.getElementById("value-email").focus();
    },
    onClickEmailList: function() {
        index.sendEmailList();
    },
    onClickLogin: function() {
        index.sendLogin(document.getElementById("value-password").value);
    },
    onClickLogout: function() {
        index.sendLogout();
    },
    onClickSignUp: function() {
        index.sendSignUp(document.getElementById("value-password").value);
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
    sendIndex: function() {
        asticode.loader.show();
        astilectron.send({name: "index"});
    },
    sendLogin: function(password) {
        asticode.loader.show();
        astilectron.send({name: "login", payload: password});
    },
    sendLogout: function(password) {
        asticode.loader.show();
        astilectron.send({name: "logout"});
    },
    sendSignUp: function(password) {
        asticode.loader.show();
        astilectron.send({name: "sign.up", payload: password});
    }
};