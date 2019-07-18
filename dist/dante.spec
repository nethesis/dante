# avoid error during caronte packaging
%define debug_package %{nil}
# do not strip caronte binary
%define __strip /bin/true

Name:		dante
Version:	0.0.10
Release:	0%{?dist}
Summary:	Single stack reports made simple

License:	GPLv3
URL:	    https://github.com/nethesis/dante	
Source0:	https://github.com/nethesis/dante/archive/master.tar.gz
# Execute ./prep-source to create Source1, Source3 and Source4
Source1:    caronte.tar.gz
Source2:    dante.sysconf
Source3:    beatrice.tar.gz
Source4:    virgilio
Source5:    virgilio.service
Source6:    dante.cron
Source7:    dante.conf

BuildRequires: systemd
BuildRequires: python

%description
Single stack reports made simple

%prep
%setup -q -n dante-master

%build
python -m compileall ciacco/lib/squidguard.py

%post
%systemd_post virgilio.service

%preun
%systemd_preun virgilio.service

%postun
%systemd_postun_with_restart virgilio.service

%install
mkdir -p %{buildroot}/usr/share/dante/caronte
mkdir -p %{buildroot}/usr/share/dante/beatrice
mkdir -p %{buildroot}/usr/share/dante/virgilio
mkdir -p %{buildroot}/usr/bin
mkdir -p %{buildroot}/etc/sysconfig/
mkdir -p %{buildroot}/etc/cron.d/
mkdir -p %{buildroot}/%{_unitdir}
mkdir -p %{buildroot}/etc/httpd/conf.d/
mkdir -p %{buildroot}/%{python_sitelib}
mv ciacco/lib/squidguardlib.py* %{buildroot}%{python_sitelib}
rm -rf ciacco/lib/
cp ciacco/ciacco %{buildroot}/%{_bindir}
cp %{SOURCE4} %{buildroot}/%{_bindir}
mv ciacco/miners %{buildroot}/usr/share/dante/
tar xvzf %{SOURCE1} -C %{buildroot}/usr/share/dante/caronte
mv %{SOURCE2}  %{buildroot}/etc/sysconfig/dante
tar xvzf %{SOURCE3} -C %{buildroot}/usr/share/dante/beatrice
cp %{SOURCE5} %{buildroot}/%{_unitdir}
cp %{SOURCE6} %{buildroot}/etc/cron.d/dante
cp %{SOURCE7} %{buildroot}/etc/httpd/conf.d/


%files
%doc README.md
%license LICENSE
%config(noreplace) /etc/sysconfig/dante
%config(noreplace) /etc/cron.d/dante
%config(noreplace) /etc/httpd/conf.d/dante.conf
%config(noreplace) /usr/share/dante/beatrice/config.js
%dir /usr/share/dante/
%dir %attr(0755, nobody, nobody) /usr/share/dante/virgilio
%{_unitdir}/virgilio.service
%{_bindir}/ciacco
%{_bindir}/virgilio
%{python_sitelib}/squidguardlib.py*
/usr/share/dante/miners/
/usr/share/dante/beatrice
/usr/share/dante/caronte


%changelog
* Wed Jul 03 2019 Giacomo Sanchietti <giacomo.sanchietti@nethesis.it> 0.0.9-0
- First alpha release
