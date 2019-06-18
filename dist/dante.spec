%define debug_package %{nil}

Name:		dante
Version:	0.0.0
Release:	0%{?dist}
Summary:	Single stack reports made simple

License:	GPLv3
URL:	    https://github.com/nethesis/dante	
Source0:	https://github.com/nethesis/dante/archive/master.tar.gz
# Execute ./prep-source to create Source1
Source1:    caronte.tar.gz
Source2:    dante.sysconf
Source3:    beatrice.tar.gz


%description
Single stack reports made simple

%prep
%setup -q -n dante-master

%package caronte
Summary: Caronte package for Dante
Requires: dante
%description caronte
Caronte creates the report preview using NodeJS and puppeteer.


%install
mkdir -p %{buildroot}/usr/share/dante/caronte
mkdir -p %{buildroot}/usr/share/dante/beatrice
mkdir -p %{buildroot}/usr/bin
mkdir -p %{buildroot}/etc/sysconfig/
cp ciacco/ciacco %{buildroot}/%{_bindir}
mv ciacco/miners %{buildroot}/usr/share/dante/
tar xvzf %{SOURCE1} -C %{buildroot}/usr/share/dante/caronte
mv %{SOURCE2}  %{buildroot}/etc/sysconfig/dante
tar xvzf %{SOURCE3} -C %{buildroot}/usr/share/dante/beatrice


%files
%doc README.md
%license LICENSE
%config /etc/sysconfig/dante
%dir /usr/share/dante/
%{_bindir}/ciacco
/usr/share/dante/miners/
/usr/share/dante/beatrice

%files caronte
/usr/share/dante/caronte


%changelog

