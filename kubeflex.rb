# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Kubeflex < Formula
  desc ""
  homepage "https://github.com/kubestellar/kubeflex"
  version "0.2.0"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/kubestellar/kubeflex/releases/download/v0.2.0/kubeflex_0.2.0_darwin_amd64.tar.gz"
      sha256 "4f8fb8ba8f247c21b99e9f0ccf5bea4cf44f927f08b53dec06a305b34529cb0c"

      def install
        bin.install "bin/kflex"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/kubestellar/kubeflex/releases/download/v0.2.0/kubeflex_0.2.0_darwin_arm64.tar.gz"
      sha256 "79ae013f646989a6db13c0a80a82c9ef37d240a76068fbbca27b7f417485c4c7"

      def install
        bin.install "bin/kflex"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/kubestellar/kubeflex/releases/download/v0.2.0/kubeflex_0.2.0_linux_arm64.tar.gz"
      sha256 "f36f40d0db03a4f42565d3571f9f214bd510db73f41aaa669621f78d83bf9900"

      def install
        bin.install "bin/kflex"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/kubestellar/kubeflex/releases/download/v0.2.0/kubeflex_0.2.0_linux_amd64.tar.gz"
      sha256 "d114abc81834ab91d304358419927b305f15a184d8683141025d3826b60d645b"

      def install
        bin.install "bin/kflex"
      end
    end
  end
end
