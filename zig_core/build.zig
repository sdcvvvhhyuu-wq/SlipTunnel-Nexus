const std = @import("std");

pub fn build(b: *std.Build) void {
    const target = b.standardTargetOptions(.{});
    const optimize = b.standardOptimizeOption(.{});

    const lib = b.addSharedLibrary(.{
        .name = "metadata_scrubber",
        .root_source_file = .{ .path = "src/metadata_scrubber.zig" },
        .target = target,
        .optimize = optimize,
    });

    // خط lib.linkLibC() به طور کامل حذف شد.
    // کد ما کاملاً Bare-Metal است و نیازی به کتابخانه استاندارد C اندروید ندارد.

    b.installArtifact(lib);
}
