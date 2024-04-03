inputs: {
  pkgs,
  config ? pkgs.config,
  lib ? pkgs.lib,
  system,
  self,
  ...
}: let
  cfg = config.services.lastfm-status;
in {
  options.services.lastfm-status = {
    enable = lib.mkEnableOption "lastfm-status";

    package = lib.mkOption {
      description = "package to use";
      default = self.packages.${system}.default;
    };

    cacheLength = lib.mkOption {
      type = lib.types.str;
      default = "1m";
      description = "how long to cache an entry for, accepts a golang time duration";
    };

    enableRatelimiting = lib.mkOption {
      type = lib.types.bool;
      default = true;
      description = "if to enable ratelimiting on the api";
    };

    port = lib.mkOption {
      type = lib.types.int;
      description = "port to run http api on";
    };
  };

  config = lib.mkIf cfg.enable {
    systemd.services.lastfm-status = {
      enable = true;
      serviceConfig = {
        DynamicUser = true;
        ProtectSystem = "full";
        ProtectHome = "yes";
        DeviceAllow = [""];
        LockPersonality = true;
        MemoryDenyWriteExecute = true;
        PrivateDevices = true;
        ProtectClock = true;
        ProtectControlGroups = true;
        ProtectHostname = true;
        ProtectKernelLogs = true;
        ProtectKernelModules = true;
        ProtectKernelTunables = true;
        ProtectProc = "invisible";
        RestrictNamespaces = true;
        RestrictRealtime = true;
        RestrictSUIDSGID = true;
        SystemCallArchitectures = "native";
        PrivateUsers = true;
        ExecStart = "${lib.getExe cfg.package} --port=${toString cfg.port} --cache-length=${cfg.cacheLength} --ratelimit=${toString cfg.enableRatelimiting}";
        Restart = "always";
      };
      wantedBy = ["default.target"];
    };
  };
}
