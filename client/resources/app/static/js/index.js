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
                case "account.added":
                    index.listenAccountAdded(message);
                    break;
                case "account.listed":
                    index.listenAccountListed(message);
                    break;
                case "account.opened":
                    index.listenAccountOpened();
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
    listenAccountAdded: function(message) {
        asticode.modaler.hide();
        asticode.notifier.success(message.payload);
        index.sendAccountList();
    },
    listenAccountListed: function(message) {
        // Init content
        let content = `<div class="index-header">
            <button class="btn btn-success" onclick="index.onClickAccountAdd()" title="Add a new account"><i class="fa fa-plus"></i></button>
            <button class="btn btn-success" onclick="index.onClickAccountList()" title="Refresh accounts list"><i class="fa fa-refresh"></i></button>
            <button class="btn btn-success" onclick="index.onClickLogout()" title="Log out"><i class="fa fa-sign-out"></i></button>
        </div>`;

        // Loop through accounts
        content += `<div class="index-list">`;
        for (let i = 0; i < message.payload.length; i++) {
            content += `<a href="` + message.payload[i].auth_url + `" target="_blank"><div class="index-item">
                ` + message.payload[i].addr + `
            </div></a>`;
        }
        content += "</div>";

        // Set content
        document.getElementById("index").innerHTML = content;
    },
    listenAccountOpened: function() {
        document.getElementById("index").innerHTML = "Bite";
    },
    listenError: function(message) {
        asticode.notifier.error(message.payload);
    },
    listenIndexed: function(message) {
        switch (message.payload) {
            case "index":
                index.sendAccountList();
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
    onClickAccountAdd: function() {
        // Build content
        let content = document.createElement("div");
        content.innerHTML = `<input type="account" placeholder="Account" id="value-account" onkeypress="if (event.keyCode === 13) document.getElementById('btn-account').click()">
        <button class="btn btn-success btn-lg" id="btn-account" onclick="index.onClickSubmitAccount()">Add</button>`;

        // Update modal
        asticode.modaler.setContent(content);
        asticode.modaler.show();
        document.getElementById("value-account").focus();
    },
    onClickAccountUnlock: function(account, auth_url) {
        // Build content
        let content = document.createElement("iframe");
        content.src = auth_url;

        // Update modal
        asticode.modaler.setContent(content);
        asticode.modaler.show();
    },
    onClickAccountList: function() {
        index.sendAccountList();
    },
    onClickAccountOpen: function(account) {
        index.sendAccountOpen(account, document.getElementById("value-password").value);
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
    onClickSubmitAccount: function() {
        index.sendAccountAdd(document.getElementById("value-account").value);
    },
    sendAccountAdd: function(account) {
        asticode.loader.show();
        astilectron.send({name: "account.add", payload: account});
    },
    sendAccountList: function() {
        asticode.loader.show();
        astilectron.send({name: "account.list"});
    },
    sendAccountOpen: function(account, password) {
        asticode.loader.show();
        astilectron.send({name: "account.open", payload: {account: account, password: password}});
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