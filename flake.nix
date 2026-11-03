{
  inputs.nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";

  outputs =
    { nixpkgs, ... }:
    let
      eachSystem =
        func:
        nixpkgs.lib.genAttrs nixpkgs.lib.systems.flakeExposed (
          system:
          let
            pkgs = nixpkgs.legacyPackages.${system};
            tools = with pkgs; [
              go
              protobuf
              protoc-gen-go
              protoc-gen-go-grpc
            ];
          in
          func pkgs tools
        );
    in
    {
      packages = eachSystem (
        pkgs: tools: {
          default = pkgs.writeShellApplication {
            name = "generate-proto";
            runtimeInputs = tools;

            text = ''
              shopt -s globstar
              protoc -I proto \
                --go_out=proto \
                --go-grpc_out=proto \
                --go_opt=paths=source_relative \
                --go-grpc_opt=paths=source_relative \
                proto/**/*.proto
            '';
          };
        }
      );

      devShells = eachSystem (
        pkgs: tools: {
          default = pkgs.mkShell {
            packages = tools;
          };
        }
      );
    };
}
