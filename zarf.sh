rm -rf zarf-cache
../zarf/build/zarf package create . --zarf-cache=zarf-cache --log-level=debug --skip-sbom

# duration=21.598682148s for a single image
# for the entire yaml
# real    1m6.582s
# user    0m12.341s
# sys     0m15.301s