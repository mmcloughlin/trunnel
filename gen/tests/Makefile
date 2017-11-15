trunnels=$(wildcard */*.trunnel)
srcs=$(addprefix c/,$(subst .trunnel,.c,$(trunnels)))

all: $(srcs)

c/%.c c/%.h: %.trunnel
	mkdir -p $(dir $@)
	python -m trunnel --target-dir $(dir $@) $<
