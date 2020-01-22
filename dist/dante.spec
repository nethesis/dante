# avoid error during caronte packaging
%define debug_package %{nil}
# do not strip caronte binary
%define __strip /bin/true

Name:		dante
Version: 0.2.1
Release: 1%{?dist}
Summary:	Single stack reports made simple

License:	GPLv3
URL:	    https://github.com/nethesis/dante	
Source0:	dante.tar.gz
# Execute ./prep-source to create Source1, Source3 and Source4
Source1:    caronte
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
%setup -q -n dante

%build
python -m compileall ciacco/lib/squidguard.py

%post
%systemd_post virgilio.service

%preun
%systemd_preun virgilio.service

%postun
%systemd_postun_with_restart virgilio.service

%install
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
cp %{SOURCE1} %{buildroot}/%{_bindir}
cp %{SOURCE4} %{buildroot}/%{_bindir}
mv ciacco/miners %{buildroot}/usr/share/dante/
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
%{_bindir}/caronte
%{_bindir}/virgilio
%{python_sitelib}/squidguardlib.py*
/usr/share/dante/miners/
/usr/share/dante/beatrice


%changelog
* Wed Jan 22 2020 Giacomo Sanchietti <giacomo.sanchietti@nethesis.it> - 0.2.1-1
- Missing users in Nextcloud report  - Bug NethServer/dev#6034

* Fri Sep 27 2019 Giacomo Sanchietti <giacomo.sanchietti@nethesis.it> - 0.2.0-1
- Weekly reports: allow the user to enter a custom time interval - nethesis/dev#5698

* Wed Sep 18 2019 Giacomo Sanchietti <giacomo.sanchietti@nethesis.it> - 0.1.4-1
- Statistics on OpenVPN connections - NethServer/dev#5827

* Tue Aug 27 2019 Giacomo Sanchietti <giacomo.sanchietti@nethesis.it> - 0.1.3-1
- Fix mail counters

* Tue Jul 23 2019 Giacomo Sanchietti <giacomo.sanchietti@nethesis.it> - 0.1.2-1
- beatrice. small fix in horizontal resize

* Mon Jul 22 2019 Giacomo Sanchietti <giacomo.sanchietti@nethesis.it> 0.1.1-1
- Add Italian translation

* Fri Jul 19 2019 Giacomo Sanchietti <giacomo.sanchietti@nethesis.it> 0.1.0-1
- First beta release

* Thu Jul 18 2019 Giacomo Sanchietti <giacomo.sanchietti@nethesis.it> 0.0.11-1
- Remove dante-caronte sub-package
- Move caronte to bin directory

* Wed Jul 03 2019 Giacomo Sanchietti <giacomo.sanchietti@nethesis.it> 0.0.9-0
- First alpha release
