# Initial Automatic Setup

This guide will walk you through the initial setup of PMail, focusing on an automatic installation process while allowing you to specify a particular binding host (IPv4 or IPv6) for your server.

## Installation Steps

First, prepare your environment and download the PMail binary. Replace `linux_arm64.zip` with the appropriate release for your system (e.g., `linux_amd64.zip`, `windows_amd64.zip`).

```bash
mkdir -p ~/pmail-setup && cd ~/pmail-setup
wget https://github.com/kelvinzer0/PMail-IPv6/releases/download/v2.8.5/linux_arm64.zip
unzip linux_arm64.zip
mv pmail_* pmail
chmod +x pmail
```

## Configuration: Binding Host

After the initial setup, you need to configure the `binding_host` in the `config.json` file. This setting allows you to manually specify the IPv4 or IPv6 address that PMail will bind to, overriding any automatic IP detection.

The `config.json` file is located in the `config` directory relative to your PMail executable (e.g., `~/pmail-setup/config/config.json`).

Open `config.json` with your preferred text editor and add or modify the `"binding_host"` entry.

**Example `config.json` snippet:**

```json
{
  "logLevel": "debug",
  "domain": "test.domain",
  "webDomain": "mail.test.domain",
  "binding_host": "0.0.0.0", // Set your desired IPv4 or IPv6 address here
  // ... other configurations
}
```

-   Set `"binding_host"` to `"0.0.0.0"` to bind to all available IPv4 interfaces.
-   Set `"binding_host"` to `"::"` to bind to all available IPv6 interfaces.
-   Set `"binding_host"` to a specific IPv4 address (e.g., `"192.168.1.100"`) or IPv6 address (e.g., `"2001:0db8::1"`) to bind only to that address.

## Running PMail

Once configured, you can start the PMail server:

```bash
./pmail
```

For background operation, you might use:

```bash
nohup ./pmail &
```

## PMail Systemctl Service

To manage PMail as a system service, you can create a systemd unit file.

1.  **Create the service file:**
    Create a file named `pmail.service` in `/etc/systemd/system/` with the following content.
    Remember to replace `/path/to/your/pmail` with the actual path where you extracted PMail.

    ```ini
    [Unit]
    Description=PMail Server
    After=network.target

    [Service]
    Type=simple
    User=root
    WorkingDirectory=/path/to/your/pmail
    ExecStart=/path/to/your/pmail/pmail
    Restart=on-failure

    [Install]
    WantedBy=multi-user.target
    ```

2.  **Reload systemd, enable, and start the service:**

    ```bash
    sudo systemctl daemon-reload
    sudo systemctl enable pmail
    sudo systemctl start pmail
    ```

3.  **Check the service status:**

    ```bash
    sudo systemctl status pmail
    ```

4.  **Stop/Restart the service:**

    ```bash
    sudo systemctl stop pmail
    sudo systemctl restart pmail
    ```
