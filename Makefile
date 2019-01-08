.PHONY: apps all clean

APPS=app_dccn_dcmgr app_dccn_taskmgr app_dccn_usermgr gateway web
RECURSIVE_MAKE = @for subdir in $(APPS); \
        do \
        echo "making in $subdir"; \
        ( cd $subdir && $(MAKE) all -f Makefile || exit 1; \
        done

RECURSIVE_CLEAN= @for subdir in $(APPS); \
        do \
        echo "cleaning in $subdir"; \
        ( cd $subdir && $(MAKE) clean -f Makefile) || exit 1; \
        done


apps:
	$(RECURSIVE_MAKE)

all: apps

clean:
	$(RECURSIVE_CLEAN)

