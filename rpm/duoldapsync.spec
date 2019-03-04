%global rev             %(git rev-parse HEAD)
%global shortrev        %(r=%{rev}; echo ${r:0:12})
%global _dwz_low_mem_die_limit 0
%define function gobuild { go build -a -ldflags "-B 0x$(head -c20 /dev/urandom|od -An -tx1|tr -d ' \n')" -v -x "$@"; }

Name:       duoldapsync
Version:    0
Release:    0.5.git%{shortrev}%{?dist}
License:    ASL 2.0
Summary:    LDAP to Duo API User and Group Syncing Daemon
Group:      System Environment/Daemons
Url:        https://github.com/bensallen/duoldapsync
Source0:    https://github.com/bensallen/%{name}/archive/%{rev}.tar.gz#/%{name}-%{rev}.tar.gz
Requires(pre):  shadow-utils

# e.g. el6 has ppc64 arch without gcc-go, so EA tag is required
ExclusiveArch:  %{?go_arches:%{go_arches}}%{!?go_arches:%{ix86} x86_64 %{arm}}
# If go_compiler is not set to 1, there is no virtual provide. Use golang instead.
BuildRequires:  %{?go_compiler:compiler(go-compiler)}%{!?go_compiler:golang}

%description
Duoldapsync is a LDAP to Duo API User and Group Syncing Daemon

%prep
%setup -q -n %{name}-%{rev}

%build
mkdir -p src/github.com/bensallen
ln -s ../../../ src/github.com/bensallen/duoldapsync

%install
export GOPATH=$(pwd):%{gopath}
# Server
%gobuild -o %{buildroot}%{_sbindir}/%{name} github.com/bensallen/duoldapsync

install -d %{buildroot}%{_unitdir}
install -m 0644 rpm/duoldapsync.service %{buildroot}%{_unitdir}/%{name}.service
install -d %{buildroot}/%{_sysconfdir}/sysconfig
install -m 0644 rpm/%{name}.sysconfig %{buildroot}/%{_sysconfdir}/sysconfig/%{name}
install -d %{buildroot}%{_sysconfdir}/duoldapsync
install -m 0644 examples/duoldapconfig.json %{buildroot}%{_sysconfdir}/%{name}/duoldapsync.json.example

%pre
getent group duoldapsync >/dev/null || groupadd -r duoldapsync
getent passwd duoldapsync >/dev/null || \
    useradd --system --gid duoldapsync --shell /sbin/nologin --home-dir %{_sysconfdir}/%{name} \
    --comment "Duoldapsync user" duoldapsync
exit 0

%post
/usr/bin/systemctl daemon-reload >/dev/null 2>&1

%preun
if [ $1 -eq 0 ] ; then
    /usr/bin/systemctl stop %{name} >/dev/null 2>&1
    /usr/bin/systemctl disable %{name} >/dev/null 2>&1
fi

%postun
if [ "$1" -ge "1" ] ; then
   /usr/bin/systemctl try-restart %{name} >/dev/null 2>&1 || :
fi

%clean
rm -rf %{buildroot}

%files
%defattr(-,root,root,-)
%{_sbindir}/%{name}
%{_unitdir}/%{name}.service
%config %{_sysconfdir}/sysconfig/%{name}
%{_sysconfdir}/%{name}/duoldapsync.json.example


%changelog
* Mon Mar 04 2019 Ben Allen <bsallen@alcf.anl.gov> - 0-0.5.gitaf02d787236b
- Bump to 0.5
* Fri Sep 07 2018 Ben Allen <bsallen@alcf.anl.gov> - 0-0.4.gita08d3a821c7a
- Bump to 0.4
* Wed Sep 05 2018 Ben Allen <bsallen@alcf.anl.gov> - 0-0.3.gita4c55289122e
- Bump to 0.3
* Wed Sep 05 2018 Ben Allen <bsallen@alcf.anl.gov> - 0-0.2.gita832051fa7b4
- Bump to 0.2
* Wed Sep 05 2018 Ben Allen <bsallen@alcf.anl.gov> - 0-0.1.gite213d7c995ca
- Initial RPM release
