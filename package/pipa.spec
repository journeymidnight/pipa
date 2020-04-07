%global debug_package %{nil}
%global __strip /bin/true

Name:           pipa
Version:        %{version}
Release:        %{release}

Summary:	pipa is a go module of process picture for yig

Group:		YIG
License:        Apache-2.0
URL:		http://github.com/journeymidnight/pipa
Source0:	%{name}_%{version}_%{release}.tar.gz
BuildRoot:	%(mktemp -ud %{_tmppath}/%{name}-%{version}-%{release}-XXXXXX)

%description


%prep
%setup -q -n %{name}_%{version}_%{release}

%install
rm -rf %{buildroot}
install -D -m 755 pipa %{buildroot}%{_bindir}/pipa
install -D -m 644 pipa.logrotate %{buildroot}/etc/logrotate.d/restore.logrotate
install -D -m 644 pipa.service %{buildroot}/usr/lib/systemd/system/yig-restore.service
install -D -m 644 pipa.toml %{buildroot}%{_sysconfdir}/yig/yig-restore.toml
install -d %{buildroot}/var/log/pipa/

%post
systemctl enable pipa

%preun

%clean
rm -rf %{buildroot}

%files
%defattr(-,root,root,-)
%config(noreplace) /etc/yig/pipa.toml
/usr/bin/pipa
/etc/logrotate.d/pipa.logrotate
%dir /var/log/pipa/
/usr/lib/systemd/system/pipa.service


%changelog