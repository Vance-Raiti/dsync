%global debug_package %{nil}

Name: dsync
Version: 0.1
Release: %{REV}%{?dist}
Summary: Synchronize profile between computers

License: GPLv2
URL: https://github.com/Vance-Raiti/dsync
Source0: %{name}-%{version}.tar.gz
Source1: dsync.service
Source2: dsync-hourly.service
Source3: dsync-hourly.timer
Source4: dsync
Source5: dsync-sync

BuildRequires: git-core
BuildRequires: systemd
BuildArch: noarch

%description
Install and enable dsync

%prep
%setup -q -n %{name}-%{version}

%build

%install
install -d -m 755 %{buildroot}%{_usr}/bin
install -d -m 755 %{buildroot}%{_unitdir}

install -m 644 %{SOURCE1} %{buildroot}%{_unitdir}/dsync.service
install -m 644 %{SOURCE2} %{buildroot}%{_unitdir}/dsync-hourly.service
install -m 644 %{SOURCE3} %{buildroot}%{_unitdir}/dsync-hourly.timer

install -m 755 %{SOURCE4} %{buildroot}%{_usr}/bin
install -m 755 %{SOURCE5} %{buildroot}%{_usr}/bin

%files
%{_unitdir}/dsync.service
%{_unitdir}/dsync-hourly.service
%{_unitdir}/dsync-hourly.timer

%{_usr}/bin/dsync
%{_usr}/bin/dsync-sync

%post
%systemd_post dsync.service dsync-hourly.service dsync-hourly.timer
systemctl --no-reload preset dsync-hourly.timer
systemctl enable --now dsync-hourly.timer
systemctl --no-reload preset dsync.service
systemctl enable --now dsync.service


%preun
%systemd_preun dsync.service dsync-hourly.service dsync-hourly.timer

%postun
%systemd_postun_with_restart dsync.service dsync-hourly.service dsync-hourly.timer

%changelog
* Wed Aug 27 2025 Vance Raiti <vraiti77@gmail.com> - 0.1-1
- Initial package
