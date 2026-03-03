{ self, testers }:
testers.nixosTest {
  name = "lastfm-status";

  nodes.machine =
    { ... }:
    {
      imports = [
        self.nixosModules.default
      ];

      services.lastfm-status.enable = true;
    };

  testScript =
    { nodes, ... }:
    ''
      start_all()
      machine.wait_for_unit("lastfm-status.service")
      machine.wait_for_open_port(${toString nodes.machine.services.lastfm-status.port})
    '';
}
